package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "kiroku/backend/docs"
	"kiroku/backend/internal/analytics"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	httpdelivery "kiroku/backend/internal/delivery/http"
	"kiroku/backend/internal/logger"
	"kiroku/backend/internal/notification"
	"kiroku/backend/internal/repository/sqlite"
	"kiroku/backend/internal/seed"
	"kiroku/backend/internal/telegram"
	"kiroku/backend/internal/tunnel"
	"kiroku/backend/internal/webpush"
)

func main() {
	devMode := config.KIROKUDev()
	logger.Init(devMode)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	if err := config.InitDirectories(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize directories")
	}

	allowedOrigins := config.AllowedCORSOrigins()

	database, err := db.Open(config.DataDir())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}
	defer database.Close()
	log.Info().Str("path", config.FullDBPath(config.DataDir())).Msg("db: connected")

	if err := auth.EnsureGenesisUser(database); err != nil {
		log.Debug().Err(err).Msg("Failed to ensure genesis user")
		log.Fatal().Msg("Failed to ensure genesis user")
	}

	// Startup warning: log (without leaking the PIN) if the genesis PIN file
	// still matches the stored user PIN hash. Only runs when DATA_DIR is set.
	if config.DataDir() != "" {
		auth.WarnIfDefaultGenesisPinInUse(database)
	}

	tgService := telegram.NewTelegramService(database, cfg)

	if err := seed.SeedCategories(database); err != nil {
		log.Fatal().Err(err).Msg("failed to seed categories")
	}

	// Repositories
	accountsRepo := sqlite.NewAccountsRepo(database)
	categoriesRepo := sqlite.NewCategoriesRepo(database)
	transactionsRepo := sqlite.NewTransactionsRepo(database)
	templatesRepo := sqlite.NewTemplatesRepo(database)
	scheduledRepo := sqlite.NewScheduledRepo(database)
	settingsRepo := sqlite.NewSettingsRepo(database)
	reportsRepo := sqlite.NewReportsRepo(database)
	sharedRepo := sqlite.NewSharedRepo(database)
	fxRepo := sqlite.NewFXRepo(database)
	labelsRepo := sqlite.NewLabelsRepo(database)
	pushRepo := sqlite.NewPushSubscriptionRepo(database)
	goalsRepo := sqlite.NewGoalsRepo(database)

	// Notification infrastructure
	notificationRepo := sqlite.NewNotificationRepo(database)
	deliveryRepo := sqlite.NewDeliveryRepo(database)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(realIP)
	r.Use(corsMiddleware(allowedOrigins, devMode))
	r.Use(logger.RequestLoggingMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/docs/swagger.json"),
	))

	// Mount shared public routes (no auth required)
	sharedHandler := httpdelivery.NewSharedHandler(sharedRepo)
	r.Mount("/shared/public", sharedHandler.PublicRoutes())

	// API V1 routes
	r.Route("/api/v1", func(r chi.Router) {
		authHandler := auth.NewHandler(database, cfg, tgService)
		r.Mount("/auth", authHandler.Routes())

		// Public /auth routes mounted above; all remaining routes require auth
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware(cfg))

			r.Mount("/shared", sharedHandler.Routes())

			accHandler := httpdelivery.NewAccountsHandler(accountsRepo)
			r.Mount("/accounts", accHandler.Routes())

			catHandler := httpdelivery.NewCategoriesHandler(categoriesRepo)
			r.Mount("/categories", catHandler.Routes())

			goalsHandler := httpdelivery.NewGoalsHandler(goalsRepo)
			r.Mount("/goals", goalsHandler.Routes())

			txnHandler := httpdelivery.NewTransactionsHandler(transactionsRepo)
			r.Mount("/transactions", txnHandler.Routes())

			tmplHandler := httpdelivery.NewTemplatesHandler(templatesRepo)
			r.Mount("/templates", tmplHandler.Routes())

			schedHandler := httpdelivery.NewScheduledHandler(scheduledRepo)
			r.Mount("/scheduled", schedHandler.Routes())

			repHandler := httpdelivery.NewReportsHandler(reportsRepo)
			r.Mount("/reports", repHandler.Routes())

			setHandler := httpdelivery.NewSettingsHandler(settingsRepo)
			r.Mount("/settings", setHandler.Routes())

			fxHandler := httpdelivery.NewFXHandler(fxRepo)
			r.Mount("/fx", fxHandler.Routes())

			// Notification service setup
			pushSender, err := webpush.NewSender()
			if err != nil {
				log.Warn().Err(err).Msg("webpush: failed to initialize sender; notifications disabled")
				pushSender = nil
			}
			notifSvc := notification.NewService(
				notificationRepo,
				deliveryRepo,
				pushRepo,
				pushSender,
				notification.ServiceConfig{
					MaxAttempts:  config.NotificationWorkerMaxAttempts(),
					RetryBase:    config.NotificationWorkerRetryBase(),
					LockTimeout:  config.NotificationWorkerLockTimeout(),
				},
			)

			// Start the delivery worker
			notifCtx, notifCancel := context.WithCancel(context.Background())
			worker := notification.NewWorker(notifSvc, notification.Config{
				Interval:     config.NotificationWorkerInterval(),
				BatchSize:    config.NotificationWorkerBatchSize(),
				MaxAttempts:  config.NotificationWorkerMaxAttempts(),
				RetryBase:    config.NotificationWorkerRetryBase(),
				LockTimeout:  config.NotificationWorkerLockTimeout(),
				Enabled:      config.NotificationWorkerEnabled(),
			}, notifCtx)
			_ = notifCtx
			defer notifCancel()
			defer worker.Stop()

			// Mount notification handler
			notificationsHandler := httpdelivery.NewNotificationsHandler(notifSvc)
			r.Mount("/notifications", notificationsHandler.Routes())

			analyticsHandler := analytics.NewHandler(database, config.FXBaseCurrency())
			r.Mount("/analytics", analyticsHandler.Routes())

			labelsHandler := httpdelivery.NewLabelsHandler(labelsRepo)
			r.Mount("/labels", labelsHandler.Routes())

			pushHandler := httpdelivery.NewPushHandler(pushRepo, cfg)
			r.Mount("/push", pushHandler.Routes())
		})
	})

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// Start FX initial fetch & scheduler
	go func() {
		var count int
		database.QueryRow("SELECT COUNT(*) FROM fx_rates").Scan(&count)
		if count == 0 {
			_ = sqlite.FetchAndStoreFX(database, config.FXBaseCurrency())
		}
	}()

	if config.TunnelEnabled() {
		urlFile := config.BackendURLPath()
		tunnelMgr := tunnel.NewManager("http://localhost"+cfg.Addr, urlFile, config.UpstashRedisRestURL(), config.UpstashRedisRestToken())
		_ = tunnelMgr.Start()
		defer tunnelMgr.Stop()
	}

	go func() {
		log.Info().Str("addr", cfg.Addr).Msg("Starting HTTP server (Clean Architecture + DDD + Repository)")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	log.Info().Msg("Server exited")
}

func realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(allowedOrigins []string, devMode bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			isAllowed := false
			if len(allowedOrigins) > 0 {
				for _, allowed := range allowedOrigins {
					if origin == allowed {
						isAllowed = true
						break
					}
				}
			} else if devMode {
				devDefaults := []string{
					"http://localhost:5173",
					"http://localhost:3000",
					"http://127.0.0.1:5173",
					"http://127.0.0.1:3000",
					"http://[::1]:5173",
					"http://[::1]:3000",
				}
				for _, allowed := range devDefaults {
					if origin == allowed {
						isAllowed = true
						break
					}
				}
			}

			if isAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

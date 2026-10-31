package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "kiroku/backend/docs"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	httpdelivery "kiroku/backend/internal/delivery/http"
	"kiroku/backend/internal/telegram"
	"kiroku/backend/internal/logger"
	"kiroku/backend/internal/repository/sqlite"
	"kiroku/backend/internal/tunnel"
)

func main() {
	devMode := os.Getenv("KIROKU_DEV") == "true"
	logger.Init(devMode)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}
	defer database.Close()

	tgService := telegram.NewTelegramService(database, cfg)

	seedCategories(database)

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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(realIP)
	r.Use(corsMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(loggingMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/docs/swagger.json"),
	))

	// Mount shared public routes
	sharedHandler := httpdelivery.NewSharedHandler(sharedRepo)
	r.Mount("/shared", sharedHandler.Routes())

	// API V1 routes
	r.Route("/api/v1", func(r chi.Router) {
		authHandler := auth.NewHandler(database, cfg, tgService)
		r.Mount("/auth", authHandler.Routes())

		accHandler := httpdelivery.NewAccountsHandler(accountsRepo)
		r.Mount("/accounts", accHandler.Routes())

		catHandler := httpdelivery.NewCategoriesHandler(categoriesRepo)
		r.Mount("/categories", catHandler.Routes())

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

		analyticsHandler := httpdelivery.NewAnalyticsHandler(sqlite.NewAnalyticsRepo(database))
		r.Mount("/analytics", analyticsHandler.Routes())

		labelsHandler := httpdelivery.NewLabelsHandler(labelsRepo)
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware(cfg))
			r.Mount("/labels", labelsHandler.Routes())
		})

		pushHandler := httpdelivery.NewPushHandler(pushRepo, cfg)
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware(cfg))
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
			_ = sqlite.FetchAndStoreFX(database, cfg.FXBaseCurrency)
		}
	}()

	if cfg.TunnelEnabled {
		urlFile := filepath.Join(cfg.DataDir, ".working", "backend-url.txt")
		tunnelMgr := tunnel.NewManager("http://localhost"+cfg.Addr, urlFile, cfg.UpstashRedisRestURL, cfg.UpstashRedisRestToken)
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

func seedCategories(database *db.DB) {
	var count int
	_ = database.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if count > 0 {
		return
	}
	seedData := []struct {
		ID       string
		ParentID *string
		Name     string
		Icon     string
		Color    string
	}{
		{"cat_food", nil, "Food", "🍽️", "#ff6b6b"},
		{"cat_transport", nil, "Transportation", "🚌", "#4ecdc4"},
		{"cat_housing", nil, "Housing", "🏠", "#ffe66d"},
		{"cat_income", nil, "Income", "💴", "#90ee90"},
		{"cat_other", nil, "Other", "📦", "#ffdac1"},
	}
	for _, c := range seedData {
		_, _ = database.Exec("INSERT INTO categories (id, parent_id, name, icon, color, sort_order) VALUES (?, ?, ?, ?, ?, 1)", c.ID, c.ParentID, c.Name, c.Icon, c.Color)
	}
}

func realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Debug().Dur("duration", time.Since(start)).Str("path", r.URL.Path).Msg("HTTP request")
	})
}

// @title Kiroku API
// @version 1.0
// @description Personal Finance Tracker API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "kiroku/backend/docs"
	"kiroku/backend/internal/accounts"
	"kiroku/backend/internal/analytics"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/backup"
	"kiroku/backend/internal/categories"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/fx"
	"kiroku/backend/internal/impexp"
	"kiroku/backend/internal/logger"
	"kiroku/backend/internal/receipts"
	"kiroku/backend/internal/reports"
	"kiroku/backend/internal/scheduled"
	"kiroku/backend/internal/settings"
	"kiroku/backend/internal/shared"
	"kiroku/backend/internal/telegram"
	"kiroku/backend/internal/templates"
	"kiroku/backend/internal/transactions"
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

	// Seed categories if empty
	seedCategories(database)

	// Initialize Telegram service
	tgService := telegram.NewTelegramService(database, cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(realIP)
	r.Use(corsMiddleware)
	r.Use(middleware.Recoverer)
	if devMode {
		r.Use(middleware.Logger)
	}

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger UI
	r.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/docs/swagger.json"),
	))

	authHandler := auth.NewHandler(database, cfg, tgService)
	authHandler.CheckAndResetPinHash()

	tgHandler := tgService.Routes()
	r.Mount("/api/v1/telegram", tgHandler)

	// Public shared access routes (no auth required)
	sharedHandler := shared.NewHandler(database)
	r.Mount("/shared", sharedHandler.Routes())

	// Public auth routes (login, logout, pin reset)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
		r.Post("/pin", authHandler.ChangePin)
		r.Post("/pin/change-request", authHandler.RequestPinChange)
		r.Post("/pin/verify", authHandler.VerifyPinChange)
		r.Post("/pin/reset-request", authHandler.RequestPinReset)
		r.Post("/pin/reset-verify", authHandler.VerifyPinReset)
		r.Post("/telegram/register-token", authHandler.GenerateRegisterToken)
		r.Get("/telegram/status", authHandler.GetTelegramStatus)
	})

	backupHandler := backup.NewHandler(database, cfg)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(cfg))

		r.Get("/api/v1/auth/me", authHandler.Me)

		accountsHandler := accounts.NewHandler(database)
		r.Mount("/api/v1/accounts", accountsHandler.Routes())

		categoriesHandler := categories.NewHandler(database)
		r.Mount("/api/v1/categories", categoriesHandler.Routes())

		transactionsHandler := transactions.NewHandler(database)
		r.Mount("/api/v1/transactions", transactionsHandler.Routes())

		templatesHandler := templates.NewHandler(database)
		r.Mount("/api/v1/templates", templatesHandler.Routes())

		scheduledHandler := scheduled.NewHandler(database)
		r.Mount("/api/v1/scheduled", scheduledHandler.Routes())

		reportsHandler := reports.NewHandler(database, cfg)
		r.Mount("/api/v1/reports", reportsHandler.Routes())

		r.Mount("/api/v1/backup", backupHandler.Routes())

		fxHandler := fx.NewHandler(database, cfg)
		r.Mount("/api/v1/fx", fxHandler.Routes())

		settingsHandler := settings.NewHandler(database)
		r.Mount("/api/v1/settings", settingsHandler.Routes())

		analyticsHandler := analytics.NewHandler(database, cfg.FXBaseCurrency)
		r.Mount("/api/v1/analytics", analyticsHandler.Routes())

		impexpHandler := impexp.NewHandler(database)
		r.Mount("/api/v1", impexpHandler.Routes())

		receiptsHandler := receipts.NewHandler(database, cfg)
		r.Mount("/api/v1/receipts", receiptsHandler.Routes())
	})

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// Start FX scheduler
	fx.StartScheduler(database, cfg.FXBaseCurrency, cfg.BackupSchedule)

	// Start backup scheduler (XLSX -> local git repo -> GitHub when configured)
	backup.StartScheduler(backupHandler.Scheduler())

	// Start Cloudflare tunnel manager & coordinate Telegram webhook
	if cfg.TunnelEnabled {
		urlFile := filepath.Join(cfg.DataDir, ".working", "backend-url.txt")
		tunnelMgr := tunnel.NewManager("http://localhost"+cfg.Addr, urlFile, cfg.UpstashRedisRestURL, cfg.UpstashRedisRestToken)
		if err := tunnelMgr.Start(); err != nil {
			log.Warn().Err(err).Msg("tunnel: failed to start")
		} else {
			go func() {
				for tunnelURL := range tunnelMgr.URLChan() {
					webhookURL := strings.TrimRight(tunnelURL, "/") + "/api/v1/telegram/webhook"
					if err := tgService.SetWebhook(webhookURL); err != nil {
						log.Warn().Err(err).Msg("telegram: failed to set webhook")
					}
				}
			}()
		}
		defer tunnelMgr.Stop()
	}

	go func() {
		log.Info().Str("addr", cfg.Addr).Msg("Starting HTTP server")
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func seedCategories(db *db.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil || count > 0 {
		return
	}

	// Japanese household taxonomy seed data
	seedData := []struct {
		ID       string
		ParentID *string
		Name     string
		Icon     string
		Color    string
		Kind     *string
		Sort     int
	}{
		// Root categories
		{"cat_food", nil, "食費", "🍽️", "#ff6b6b", nil, 1},
		{"cat_transport", nil, "交通費", "🚌", "#4ecdc4", nil, 2},
		{"cat_housing", nil, "住居費", "🏠", "#ffe66d", nil, 3},
		{"cat_utilities", nil, "光熱費", "💡", "#95e1d3", nil, 4},
		{"cat_communication", nil, "通信費", "📱", "#f38181", nil, 5},
		{"cat_insurance", nil, "保険", "🛡️", "#aa96da", nil, 6},
		{"cat_medical", nil, "医療費", "🏥", "#fcbad3", nil, 7},
		{"cat_education", nil, "教育費", "📚", "#a8d8ea", nil, 8},
		{"cat_entertainment", nil, "娯楽費", "🎮", "#f8b500", nil, 9},
		{"cat_clothing", nil, "被服費", "👕", "#ffb3ba", nil, 10},
		{"cat_savings", nil, "貯蓄", "💰", "#b5ead7", nil, 11},
		{"cat_investment", nil, "投資", "📈", "#c7ceea", nil, 12},
		{"cat_income", nil, "収入", "💴", "#90ee90", nil, 13},
		{"cat_other", nil, "その他", "📦", "#ffdac1", nil, 14},

		// Food subcategories
		{"cat_food_groceries", strPtr("cat_food"), "食材", "🥬", "#ff8e8e", strPtr("must"), 1},
		{"cat_food_dining", strPtr("cat_food"), "外食", "🍜", "#ff8e8e", strPtr("have"), 2},
		{"cat_food_cafe", strPtr("cat_food"), "カフェ", "☕", "#ff8e8e", strPtr("want"), 3},
		{"cat_food_snacks", strPtr("cat_food"), "お菓子", "🍪", "#ff8e8e", strPtr("want"), 4},
		{"cat_food_drinks", strPtr("cat_food"), "飲料", "🥤", "#ff8e8e", strPtr("have"), 5},

		// Transport subcategories
		{"cat_transport_train", strPtr("cat_transport"), "電車", "🚃", "#70e0d8", strPtr("must"), 1},
		{"cat_transport_bus", strPtr("cat_transport"), "バス", "🚌", "#70e0d8", strPtr("must"), 2},
		{"cat_transport_taxi", strPtr("cat_transport"), "タクシー", "🚕", "#70e0d8", strPtr("want"), 3},
		{"cat_transport_gas", strPtr("cat_transport"), "ガソリン", "⛽", "#70e0d8", strPtr("must"), 4},
		{"cat_transport_toll", strPtr("cat_transport"), "高速代", "🛣️", "#70e0d8", strPtr("must"), 5},

		// Housing subcategories
		{"cat_housing_rent", strPtr("cat_housing"), "家賃", "🏠", "#fff08f", strPtr("must"), 1},
		{"cat_housing_electric", strPtr("cat_housing"), "電気代", "⚡", "#fff08f", strPtr("must"), 2},
		{"cat_housing_gas", strPtr("cat_housing"), "ガス代", "🔥", "#fff08f", strPtr("must"), 3},
		{"cat_housing_water", strPtr("cat_housing"), "水道代", "💧", "#fff08f", strPtr("must"), 4},
		{"cat_housing_maintenance", strPtr("cat_housing"), "修繕費", "🔧", "#fff08f", strPtr("have"), 5},

		// Utilities subcategories
		{"cat_util_internet", strPtr("cat_utilities"), "インターネット", "🌐", "#b8ebe3", strPtr("must"), 1},
		{"cat_util_phone", strPtr("cat_utilities"), "携帯電話", "📱", "#b8ebe3", strPtr("must"), 2},
		{"cat_util_subscriptions", strPtr("cat_utilities"), "サブスク", "📺", "#b8ebe3", strPtr("want"), 3},

		// Income subcategories
		{"cat_income_salary", strPtr("cat_income"), "給与", "💼", "#a8f0a8", nil, 1},
		{"cat_income_bonus", strPtr("cat_income"), "ボーナス", "🎁", "#a8f0a8", nil, 2},
		{"cat_income_side", strPtr("cat_income"), "副業", "💻", "#a8f0a8", nil, 3},
		{"cat_income_investment", strPtr("cat_income"), "投資収益", "📊", "#a8f0a8", nil, 4},
		{"cat_income_other", strPtr("cat_income"), "その他収入", "💰", "#a8f0a8", nil, 5},
	}

	for _, cat := range seedData {
		_, err := db.Exec(
			`INSERT INTO categories (id, parent_id, name, icon, color, kind, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			cat.ID, cat.ParentID, cat.Name, cat.Icon, cat.Color, cat.Kind, cat.Sort,
		)
		if err != nil {
			log.Error().Err(err).Str("category", cat.Name).Msg("seed categories")
		}
	}
	log.Info().Int("count", len(seedData)).Msg("Seeded categories")
}

func strPtr(s string) *string {
	return &s
}

func realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Safe passthrough without trusting untrusted proxy headers unless behind known proxy
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

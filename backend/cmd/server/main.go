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

		r.Use(auth.AuthMiddleware(cfg))

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
		r.Mount("/labels", labelsHandler.Routes())

		pushHandler := httpdelivery.NewPushHandler(pushRepo, cfg)
		r.Mount("/push", pushHandler.Routes())
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
		Kind     string
	}{
		// Root categories (parent_id = nil)
		{"cat_food_drink", nil, "Food & Drink", "🍽️", "#EF4444", "have"},
		{"cat_shopping", nil, "Shopping", "🛍️", "#38BDF8", "want"},
		{"cat_investment", nil, "Investment", "📊", "#EC4899", "want"},
		{"cat_income", nil, "Income", "💰", "#FACC15", "want"},
		{"cat_financial_expenses", nil, "Financial Expenses", "🧾", "#0EA5E9", "must"},
		{"cat_life_entertainment", nil, "Life & Entertainment", "✨", "#4ADE80", "want"},
		{"cat_communication_pc", nil, "Communication & PC", "💻", "#8195F5", "have"},
		{"cat_housing", nil, "Housing", "🏠", "#F59E0B", "must"},
		{"cat_vehicle", nil, "Vehicle", "🚗", "#C052DC", "have"},
		{"cat_other", nil, "Other", "📦", "#A1A1A1", "want"},

		// Food & Drink children
		{"cat_bar_cafe", strptr("cat_food_drink"), "Bar, Cafe", "🍷", "#EF4444", "want"},
		{"cat_groceries", strptr("cat_food_drink"), "Groceries", "🥦", "#EF4444", "have"},
		{"cat_restaurant_fastfood", strptr("cat_food_drink"), "Restaurant, Fast-Food", "🍕", "#EF4444", "have"},

		// Shopping children
		{"cat_clothes_shoes", strptr("cat_shopping"), "Clothes, & Shoes", "👕", "#38BDF8", "have"},
		{"cat_drug_store", strptr("cat_shopping"), "Drug-Store, Chenist", "💊", "#38BDF8", "have"},
		{"cat_electronics_accessories", strptr("cat_shopping"), "Electronics, Accessories", "💻", "#38BDF8", "want"},
		{"cat_free_time", strptr("cat_shopping"), "Free Time", "🏖️", "#38BDF8", "want"},
		{"cat_gifts_joy", strptr("cat_shopping"), "Gifts & Joy", "🎁", "#38BDF8", "want"},
		{"cat_healthy_beauty", strptr("cat_shopping"), "Healthy & Beauty", "💄", "#38BDF8", "want"},
		{"cat_home_garden", strptr("cat_shopping"), "Home, Garden", "🏡", "#38BDF8", "want"},
		{"cat_jewels_accessories", strptr("cat_shopping"), "Jewels & Accessories", "💎", "#38BDF8", "want"},
		{"cat_kids", strptr("cat_shopping"), "Kids", "👶", "#38BDF8", "have"},
		{"cat_pets_animal", strptr("cat_shopping"), "Pets, Animal", "🐾", "#38BDF8", "want"},

		// Investment children
		{"cat_realty", strptr("cat_investment"), "Realty", "🏢", "#EC4899", "want"},
		{"cat_savings", strptr("cat_investment"), "Savings", "🐖", "#EC4899", "want"},
		{"cat_vehicles_chattels", strptr("cat_investment"), "Vehicles, Chattels", "🚗", "#EC4899", "want"},

		// Income children
		{"cat_rental_income", strptr("cat_income"), "Rental Income", "🏠", "#FACC15", "want"},
		{"cat_sale", strptr("cat_income"), "Sale", "🏷️", "#FACC15", "want"},
		{"cat_wage_invoices", strptr("cat_income"), "Wage, Invoices", "📃", "#FACC15", "want"},
		{"cat_dividends_interest", strptr("cat_income"), "Dividends, Interest", "📈", "#FACC15", "want"},
		{"cat_checks_coupons", strptr("cat_income"), "Checks, Coupons", "🎫", "#FACC15", "want"},
		{"cat_gifts_income", strptr("cat_income"), "Gifts", "🎁", "#FACC15", "want"},
		{"cat_refund_tax", strptr("cat_income"), "Refund(Tax)", "💰", "#FACC15", "must"},
		{"cat_child_support_income", strptr("cat_income"), "Child Support", "👶", "#FACC15", "want"},
		{"cat_lottery_gambling_income", strptr("cat_income"), "Lottery, Gambling", "🎰", "#FACC15", "want"},

		// Financial Expenses children
		{"cat_taxes", strptr("cat_financial_expenses"), "Taxes", "🧾", "#0EA5E9", "must"},
		{"cat_child_support_expense", strptr("cat_financial_expenses"), "Child Support", "👶", "#0EA5E9", "must"},
		{"cat_fines", strptr("cat_financial_expenses"), "Fines", "⚖️", "#0EA5E9", "must"},
		{"cat_loan_interest", strptr("cat_financial_expenses"), "Loan, Interest", "🏛️", "#0EA5E9", "must"},
		{"cat_insurance_financial", strptr("cat_financial_expenses"), "Insurances", "🛡️", "#0EA5E9", "have"},
		{"cat_advisory", strptr("cat_financial_expenses"), "Advisory", "❓", "#0EA5E9", "want"},
		{"cat_charges_fees", strptr("cat_financial_expenses"), "Charges, Fees", "📋", "#0EA5E9", "must"},

		// Life & Entertainment children
		{"cat_tv_streaming", strptr("cat_life_entertainment"), "TV, Streaming", "📺", "#4ADE80", "want"},
		{"cat_active_fitness", strptr("cat_life_entertainment"), "Active Sport & Fitness", "💪", "#4ADE80", "want"},
		{"cat_alcohol_tobacco", strptr("cat_life_entertainment"), "Alcohol & Tobacco", "🍷", "#4ADE80", "want"},
		{"cat_books_subscriptions", strptr("cat_life_entertainment"), "Books, Audio, Subscriptions", "📚", "#4ADE80", "want"},
		{"cat_charity_gifts", strptr("cat_life_entertainment"), "Charity & Gifts", "🎁", "#4ADE80", "want"},
		{"cat_culture_sport_events", strptr("cat_life_entertainment"), "Culture, Sport Events", "🎟️", "#4ADE80", "want"},
		{"cat_education_development", strptr("cat_life_entertainment"), "Education, Development", "🎓", "#4ADE80", "have"},
		{"cat_health_care_doctor", strptr("cat_life_entertainment"), "Health Care, Doctor", "🏥", "#4ADE80", "have"},
		{"cat_hobbies", strptr("cat_life_entertainment"), "Hobbies", "🎨", "#4ADE80", "want"},
		{"cat_holiday_trips_events", strptr("cat_life_entertainment"), "Holiday, Trips, Events", "🌴", "#4ADE80", "want"},
		{"cat_lottery_gambling_life", strptr("cat_life_entertainment"), "Lottery, Gambling", "🎰", "#4ADE80", "want"},

		// Communication & PC children
		{"cat_internet", strptr("cat_communication_pc"), "Internet", "🌐", "#8195F5", "have"},
		{"cat_phone_cell_phone", strptr("cat_communication_pc"), "Phone, Cell Phone", "📱", "#8195F5", "have"},
		{"cat_postal_services", strptr("cat_communication_pc"), "Postal Services", "📬", "#8195F5", "have"},
		{"cat_software_apps_games", strptr("cat_communication_pc"), "Software, Apps, Games", "🎮", "#8195F5", "want"},

		// Housing children
		{"cat_energy_utilities", strptr("cat_housing"), "Energy, Utilities", "💡", "#F59E0B", "must"},
		{"cat_maintenances_repairs", strptr("cat_housing"), "Maintenances, Repairs", "🔧", "#F59E0B", "have"},
		{"cat_mortgage", strptr("cat_housing"), "Mortgage", "🏦", "#F59E0B", "must"},
		{"cat_property_insurance", strptr("cat_housing"), "Property Insurance", "🏠", "#F59E0B", "have"},
		{"cat_rent", strptr("cat_housing"), "Rent", "🏠", "#F59E0B", "must"},
		{"cat_services_housing", strptr("cat_housing"), "Services", "🛠️", "#F59E0B", "must"},

		// Vehicle children
		{"cat_fuel", strptr("cat_vehicle"), "Fuel", "⛽", "#C052DC", "have"},
		{"cat_leasing", strptr("cat_vehicle"), "Leasing", "💳", "#C052DC", "must"},
		{"cat_parking", strptr("cat_vehicle"), "Parking", "🅿️", "#C052DC", "want"},
		{"cat_rentals_vehicle", strptr("cat_vehicle"), "Rentals", "🗝️", "#C052DC", "want"},
		{"cat_vehicle_insurance", strptr("cat_vehicle"), "Vehicle Insurance", "🆔", "#C052DC", "have"},
		{"cat_vehicle_maintenance", strptr("cat_vehicle"), "Vehicle Maintenance", "🔧", "#C052DC", "have"},
	}
	for _, c := range seedData {
		var parentSQL string
		if c.ParentID != nil {
			parentSQL = *c.ParentID
		}
		_, _ = database.Exec("INSERT INTO categories (id, parent_id, name, icon, color, kind, sort_order) VALUES (?, ?, ?, ?, ?, ?, 1)", c.ID, parentSQL, c.Name, c.Icon, c.Color, c.Kind)
	}
}

func strptr(s string) *string { return &s }

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

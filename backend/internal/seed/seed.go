package seed

import (
	"context"
	"fmt"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/logger"
)

func SeedCategories(db *db.DB) error {
	ctx := context.Background()

	var count int
	_ = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories").Scan(&count)
	if count > 0 {
		logger.GetLogger().Info().Msg("categories already seeded, skipping")
		return nil
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
		_, err := db.ExecContext(ctx, "INSERT INTO categories (id, parent_id, name, icon, color, kind, sort_order) VALUES (?, ?, ?, ?, ?, ?, 1)", c.ID, c.ParentID, c.Name, c.Icon, c.Color, c.Kind)
		if err != nil {
			logger.GetLogger().Error().Err(err).Str("id", c.ID).Msg("seed category insert error")
			return fmt.Errorf("seed category %s: %w", c.ID, err)
		}
		logger.GetLogger().Debug().Str("id", c.ID).Msg("seeded category")
	}

	logger.GetLogger().Info().Msg("seeded categories")
	return nil
}

func strptr(s string) *string { return &s }
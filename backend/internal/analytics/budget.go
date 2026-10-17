package analytics

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetBudgetSummary(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateTo = time.Now().Format("2006-01-02")
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT c.id, c.name, c.icon, c.color, c.budget_amount, c.budget_period,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as spent
		FROM categories c
		LEFT JOIN transactions t ON t.category_id = c.id AND t.deleted_at IS NULL AND t.date >= ? AND t.date <= ?
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE c.user_id = ? AND c.budget_amount > 0 AND c.budget_period != ''
		GROUP BY c.id, c.name, c.icon, c.color, c.budget_amount, c.budget_period
	`
	rows, err := h.db.Query(query, dateFrom, dateTo, baseCurrency, userID)
	if err != nil {
		log.Error().Err(err).Msg("analytics: budget categories")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var budgets []models.BudgetItem
	var totalBudget, totalSpent float64

	for rows.Next() {
		var item models.BudgetItem
		var spent float64
		if err := rows.Scan(&item.CategoryID, &item.CategoryName, &item.Icon, &item.Color, &item.BudgetAmount, &item.BudgetPeriod, &spent); err != nil {
			log.Error().Err(err).Msg("analytics: budget scan")
			continue
		}
		item.Spent = spent
		item.Remaining = item.BudgetAmount - spent
		if item.BudgetAmount > 0 {
			item.ProgressPct = (spent / item.BudgetAmount) * 100
		}
		if item.Remaining < 0 {
			item.Status = "over"
		} else if item.ProgressPct >= 80 {
			item.Status = "warning"
		} else {
			item.Status = "ok"
		}
		budgets = append(budgets, item)
		totalBudget += item.BudgetAmount
		totalSpent += spent
	}

	writeJSON(w, map[string]interface{}{
		"budgets": budgets,
		"summary": map[string]interface{}{
			"total_budget": totalBudget,
			"total_spent":  totalSpent,
			"total_remain": totalBudget - totalSpent,
			"overall_pct":  safeDiv(totalSpent, totalBudget) * 100,
		},
	})
}

func (h *Handler) GetBudgetVsActual(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	categoryID := r.URL.Query().Get("category_id")

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateTo = time.Now().Format("2006-01-02")
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT c.id, c.name, c.icon, c.color, c.budget_amount, c.budget_period
		FROM categories c
		WHERE c.user_id = ? AND c.budget_amount > 0 AND c.budget_period != ''
	`
	args := []interface{}{userID}
	if categoryID != "" {
		query += " AND c.id = ?"
		args = append(args, categoryID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.BudgetVsActual
	for rows.Next() {
		var item models.BudgetVsActual
		rows.Scan(&item.CategoryID, &item.CategoryName, &item.Icon, &item.Color, &item.BudgetAmount, &item.BudgetPeriod)

		var spent float64
		h.db.QueryRow(`
			SELECT COALESCE(SUM(t.amount * COALESCE(r.rate, 1)), 0)
			FROM transactions t
			LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
			WHERE t.user_id = ? AND t.category_id = ? AND t.type = 'expense'
				AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		`, baseCurrency, userID, item.CategoryID, dateFrom, dateTo).Scan(&spent)

		item.ActualSpent = spent
		item.Variance = spent - item.BudgetAmount
		item.VariancePct = safeDiv(item.Variance, item.BudgetAmount) * 100
		results = append(results, item)
	}

	writeJSON(w, map[string]interface{}{"budgets_vs_actual": results})
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

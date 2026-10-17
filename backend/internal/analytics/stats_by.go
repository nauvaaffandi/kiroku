package analytics

import (
	"net/http"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetStatsByCurrency(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT a.currency,
			COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as expense,
			COUNT(*) as count
		FROM transactions t
		JOIN accounts a ON t.from_account_id = a.id
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE t.user_id = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		GROUP BY a.currency
	`
	rows, err := h.db.Query(query, baseCurrency, userID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var currencies []models.CurrencyStats
	for rows.Next() {
		var cs models.CurrencyStats
		rows.Scan(&cs.Currency, &cs.Income, &cs.Expense, &cs.Count)
		cs.Net = cs.Income - cs.Expense
		currencies = append(currencies, cs)
	}

	writeJSON(w, map[string]interface{}{
		"currencies":    currencies,
		"base_currency": baseCurrency,
	})
}

func (h *Handler) GetStatsByAccount(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT a.id, a.name, a.currency,
			COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as expense,
			COUNT(*) as count
		FROM transactions t
		JOIN accounts a ON (t.from_account_id = a.id OR t.to_account_id = a.id)
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE t.user_id = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL AND a.archived_at IS NULL
		GROUP BY a.id, a.name, a.currency
		ORDER BY income + expense DESC
	`
	rows, err := h.db.Query(query, baseCurrency, userID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []models.AccountStats
	for rows.Next() {
		var as models.AccountStats
		rows.Scan(&as.AccountID, &as.AccountName, &as.Currency, &as.Income, &as.Expense, &as.Count)
		as.Net = as.Income - as.Expense
		accounts = append(accounts, as)
	}

	writeJSON(w, map[string]interface{}{
		"accounts":      accounts,
		"base_currency": baseCurrency,
	})
}

package analytics

import (
	"net/http"
	"time"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetCashFlowStats(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	interval := r.URL.Query().Get("interval")

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}
	if interval == "" {
		interval = "month"
	}

	start, _ := time.Parse("2006-01-02", dateFrom)
	end, _ := time.Parse("2006-01-02", dateTo)

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	var series []models.CashFlowSeriesPoint
	currStart := start

	for currStart.Before(end) {
		var currEnd time.Time
		var periodStr string
		if interval == "day" {
			currEnd = currStart.AddDate(0, 0, 1)
			periodStr = currStart.Format("2006-01-02")
		} else if interval == "week" {
			currEnd = currStart.AddDate(0, 0, 7)
			periodStr = currStart.Format("2006-01-02")
		} else {
			currEnd = currStart.AddDate(0, 1, 0)
			periodStr = currStart.Format("2006-01")
		}
		if currEnd.After(end) {
			currEnd = end
		}

		var income, expense float64
		h.db.QueryRow(`
			SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount * COALESCE(r.rate, 1) ELSE 0 END), 0),
				COALESCE(SUM(CASE WHEN type = 'expense' THEN amount * COALESCE(r.rate, 1) ELSE 0 END), 0)
			FROM transactions t
			LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
			WHERE t.user_id = ? AND t.date >= ? AND t.date < ? AND t.deleted_at IS NULL
		`, baseCurrency, userID, currStart.Format("2006-01-02"), currEnd.Format("2006-01-02")).Scan(&income, &expense)

		series = append(series, models.CashFlowSeriesPoint{
			Period:  periodStr,
			Income:  income,
			Expense: expense,
			Net:     income - expense,
		})

		currStart = currEnd
	}

	writeJSON(w, map[string]interface{}{
		"series":   series,
		"interval": interval,
	})
}

func (h *Handler) GetNetWorth(w http.ResponseWriter, r *http.Request) {
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

	rows, err := h.db.Query(`
		SELECT a.id, a.name, a.currency, COALESCE(a.balance, 0)
		FROM accounts a
		WHERE a.user_id = ? AND a.archived_at IS NULL
	`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []models.NetWorthAccount
	var totalAssets, totalLiabilities float64
	for rows.Next() {
		var acc models.NetWorthAccount
		var bal float64
		rows.Scan(&acc.AccountID, &acc.Name, &acc.Currency, &bal)

		if acc.Currency != baseCurrency {
			rate := h.getFxRate(acc.Currency, baseCurrency)
			bal *= rate
			acc.Currency = baseCurrency
		}
		acc.Balance = bal
		acc.Type = "asset"
		if bal < 0 {
			acc.Type = "liability"
			totalLiabilities += -bal
		} else {
			totalAssets += bal
		}
		accounts = append(accounts, acc)
	}

	writeJSON(w, map[string]interface{}{
		"accounts": accounts,
		"summary": map[string]interface{}{
			"total_assets":      totalAssets,
			"total_liabilities": totalLiabilities,
			"net_worth":         totalAssets - totalLiabilities,
			"currency":          baseCurrency,
		},
		"date_from": dateFrom,
		"date_to":   dateTo,
	})
}

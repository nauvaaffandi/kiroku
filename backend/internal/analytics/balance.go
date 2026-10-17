package analytics

import (
	"net/http"
	"time"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetBalanceTrend(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	interval := r.URL.Query().Get("interval")
	accountID := r.URL.Query().Get("account_id")

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateTo = time.Now().Format("2006-01-02")
	}
	if interval == "" {
		interval = "month"
	}

	start, _ := time.Parse("2006-01-02", dateFrom)
	end, _ := time.Parse("2006-01-02", dateTo)

	var series []models.BalanceTrendPoint
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

		opening, closing, net := h.calculateBalanceForRange(userID, currStart.Format("2006-01-02"), currEnd.Format("2006-01-02"), accountID)

		series = append(series, models.BalanceTrendPoint{
			Period:         periodStr,
			OpeningBalance: opening,
			ClosingBalance: closing,
			NetChange:      net,
			Currency:       h.cfg.FXBaseCurrency,
		})

		currStart = currEnd
	}

	writeJSON(w, map[string]interface{}{"series": series})
}

func (h *Handler) GetBalanceByAccounts(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		dateFrom = "2024-01-01"
		dateTo = time.Now().Format("2006-01-02")
	}

	rows, err := h.db.Query(`SELECT id, name, currency, color, COALESCE(balance, 0) as balance FROM accounts WHERE user_id = ? AND archived_at IS NULL`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accounts []models.AccountBalanceSummary
	totalByCurrency := make(map[string]float64)

	for rows.Next() {
		var acc models.AccountBalanceSummary
		var baseBalance float64
		rows.Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Color, &baseBalance)

		opening, closing, net := h.calculateBalanceForAccountRange(userID, acc.ID, dateFrom, dateTo)
		acc.OpeningBalance = opening
		acc.ClosingBalance = closing
		acc.NetChange = net

		accounts = append(accounts, acc)
		totalByCurrency[acc.Currency] += closing
	}

	writeJSON(w, map[string]interface{}{
		"accounts":          accounts,
		"total_by_currency": totalByCurrency,
	})
}

func (h *Handler) GetBalanceByCurrencies(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if dateFrom == "" || dateTo == "" {
		dateFrom = "2024-01-01"
		dateTo = time.Now().Format("2006-01-02")
	}

	rows, err := h.db.Query(`SELECT a.currency, SUM(COALESCE(a.balance, 0)) as total_bal FROM accounts a WHERE a.user_id = ? AND a.archived_at IS NULL GROUP BY a.currency`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var currencies []models.CurrencyBalanceSummary
	for rows.Next() {
		var curr string
		var totalBal float64
		rows.Scan(&curr, &totalBal)

		opening, closing, net := h.calculateBalanceForCurrencyRange(userID, curr, dateFrom, dateTo)
		currencies = append(currencies, models.CurrencyBalanceSummary{
			Currency:       curr,
			OpeningBalance: opening,
			ClosingBalance: closing,
			NetChange:      net,
			FXRateToBase:   1.0,
			ConvertedTotal: closing,
		})
	}

	writeJSON(w, map[string]interface{}{"currencies": currencies})
}

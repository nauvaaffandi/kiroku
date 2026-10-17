package analytics

import (
	"database/sql"
	"net/http"
	"time"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetPlannedByCategory(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	txnType := r.URL.Query().Get("type")

	if txnType == "" {
		txnType = "expense"
	}

	query := `
		SELECT c.id, c.name, c.icon, c.color, COALESCE(SUM(s.amount), 0) as total
		FROM scheduled_transactions s
		JOIN categories c ON s.category_id = c.id
		WHERE s.user_id = ? AND s.type = ? AND s.enabled = 1
		GROUP BY c.id, c.name, c.icon, c.color
	`
	rows, err := h.db.Query(query, userID, txnType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.CategoryPlanned
	var grandTotal float64
	for rows.Next() {
		var cp models.CategoryPlanned
		rows.Scan(&cp.CategoryID, &cp.CategoryName, &cp.Icon, &cp.Color, &cp.TotalAmount)
		categories = append(categories, cp)
		grandTotal += cp.TotalAmount
	}

	writeJSON(w, map[string]interface{}{
		"categories":  categories,
		"grand_total": grandTotal,
	})
}

func (h *Handler) GetPlannedByTimeline(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	interval := r.URL.Query().Get("interval")

	if dateFrom == "" || dateTo == "" {
		dateFrom = time.Now().Format("2006-01-02")
		dateTo = time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	}
	if interval == "" {
		interval = "month"
	}

	start, _ := time.Parse("2006-01-02", dateFrom)
	end, _ := time.Parse("2006-01-02", dateTo)

	rows, err := h.db.Query(`SELECT id, amount, type, rrule FROM scheduled_transactions WHERE user_id = ? AND enabled = 1`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type schedItem struct {
		id     int64
		amount float64
		tType  string
		rrule  sql.NullString
	}
	var items []schedItem
	for rows.Next() {
		var si schedItem
		rows.Scan(&si.id, &si.amount, &si.tType, &si.rrule)
		items = append(items, si)
	}

	var timeline []models.TimelinePoint
	curr := start
	for curr.Before(end) {
		var next time.Time
		var pStr string
		if interval == "day" {
			next = curr.AddDate(0, 0, 1)
			pStr = curr.Format("2006-01-02")
		} else if interval == "week" {
			next = curr.AddDate(0, 0, 7)
			pStr = curr.Format("2006-01-02")
		} else {
			next = curr.AddDate(0, 1, 0)
			pStr = curr.Format("2006-01")
		}
		if next.After(end) {
			next = end
		}

		var income, expense float64
		for _, item := range items {
			dates := h.expandRRULE(item.rrule.String, curr, next)
			count := float64(len(dates))
			if count > 0 {
				if item.tType == "income" {
					income += item.amount * count
				} else {
					expense += item.amount * count
				}
			}
		}

		timeline = append(timeline, models.TimelinePoint{
			Period:  pStr,
			Income:  income,
			Expense: expense,
			Net:     income - expense,
		})

		curr = next
	}

	writeJSON(w, map[string]interface{}{"timeline": timeline})
}

func (h *Handler) GetBalanceForecast(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateTo := r.URL.Query().Get("date_to")
	if dateTo == "" {
		dateTo = time.Now().AddDate(0, 6, 0).Format("2006-01-02")
	}

	start := time.Now()
	end, _ := time.Parse("2006-01-02", dateTo)

	accounts, err := h.getAccountBalances(userID, "", h.cfg.FXBaseCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var currentBalance float64
	for _, acc := range accounts {
		currentBalance += acc.Balance
	}

	scheduled, err := h.getScheduledTransactions(userID, start, end, "", h.cfg.FXBaseCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dayMap := make(map[string]float64)
	for _, item := range scheduled {
		if item.Type == "income" {
			dayMap[item.Date] += item.Amount
		} else {
			dayMap[item.Date] -= item.Amount
		}
	}

	var forecast []models.ForecastPoint
	curr := start
	running := currentBalance
	for !curr.After(end) {
		dStr := curr.Format("2006-01-02")
		delta := dayMap[dStr]
		running += delta

		forecast = append(forecast, models.ForecastPoint{
			Date:             dStr,
			ProjectedBalance: running,
			Delta:            delta,
		})
		curr = curr.AddDate(0, 0, 1)
	}

	writeJSON(w, map[string]interface{}{
		"current_balance": currentBalance,
		"forecast":        forecast,
	})
}

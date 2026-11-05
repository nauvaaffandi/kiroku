package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
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

	rows, err := h.db.Query(`SELECT id, rrule, transaction_template_json, enabled FROM scheduled WHERE enabled = 1`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := make([]models.CategoryPlanned, 0)
	var grandTotal float64
	type catAgg struct {
		id, name, icon, color string
		total                 float64
	}
	aggMap := make(map[string]*catAgg)
	for rows.Next() {
		var sid, rrule, tmplJSON string
		if err := rows.Scan(&sid, &rrule, &tmplJSON); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var tmpl map[string]interface{}
		if err := json.Unmarshal([]byte(tmplJSON), &tmpl); err != nil {
			continue
		}
		tType, _ := tmpl["type"].(string)
		if tType != txnType {
			continue
		}
		amt := 0.0
		if a, ok := tmpl["amount"].(float64); ok {
			amt = a
		}
		cid, _ := tmpl["category_id"].(string)
		if cid == "" {
			continue
		}
		var catName, catIcon, catColor string
		h.db.QueryRow(`SELECT name, icon, color FROM categories WHERE id = ?`, cid).Scan(&catName, &catIcon, &catColor)
		if agg, ok := aggMap[cid]; ok {
			agg.total += amt
		} else {
			aggMap[cid] = &catAgg{id: cid, name: catName, icon: catIcon, color: catColor, total: amt}
		}
		grandTotal += amt
	}
	for _, agg := range aggMap {
		categories = append(categories, models.CategoryPlanned{
			CategoryID:   agg.id,
			CategoryName: agg.name,
			Icon:         agg.icon,
			Color:        agg.color,
			TotalAmount:  agg.total,
		})
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

	rows, err := h.db.Query(`SELECT id, rrule, transaction_template_json, enabled FROM scheduled WHERE enabled = 1`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type schedItem struct {
		id       string
		rrule    string
		tmplJSON string
	}
	var items []schedItem
	for rows.Next() {
		var si schedItem
		if err := rows.Scan(&si.id, &si.rrule, &si.tmplJSON); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, si)
	}

	timeline := make([]models.TimelinePoint, 0)
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
			var tmpl map[string]interface{}
			if err := json.Unmarshal([]byte(item.tmplJSON), &tmpl); err != nil {
				continue
			}
			amt := 0.0
			if a, ok := tmpl["amount"].(float64); ok {
				amt = a
			}
			tType, _ := tmpl["type"].(string)
			dates := h.expandRRULE(item.rrule, curr, next)
			count := float64(len(dates))
			if count > 0 {
				if tType == "income" {
					income += amt * count
				} else {
					expense += amt * count
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

	accounts, err := h.getAccountBalances(userID, "", config.FXBaseCurrency())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var currentBalance float64
	for _, acc := range accounts {
		currentBalance += acc.Balance
	}

	scheduled, err := h.getScheduledTransactions(userID, start, end, "", config.FXBaseCurrency())
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

	forecast := make([]models.ForecastPoint, 0)
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

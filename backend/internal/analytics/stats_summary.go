package analytics

import (
	"net/http"
	"strconv"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetStatsSummary(w http.ResponseWriter, r *http.Request) {
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

	stats := h.getPeriodStats(userID, dateFrom, dateTo)
	days := calculateDays(dateFrom, dateTo)

	writeJSON(w, map[string]interface{}{
		"date_from":         dateFrom,
		"date_to":           dateTo,
		"days":              days,
		"income":            stats["income"],
		"expense":           stats["expense"],
		"net":               stats["net"],
		"txn_count":         int(stats["count"]),
		"avg_daily_income":  stats["income"] / days,
		"avg_daily_expense": stats["expense"] / days,
	})
}

func (h *Handler) GetStatsByCategory(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	txnType := r.URL.Query().Get("type")
	limitStr := r.URL.Query().Get("limit")

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}
	if txnType == "" {
		txnType = "expense"
	}
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := `
		SELECT c.id, c.name, c.icon, c.color,
			COALESCE(SUM(t.amount * COALESCE(r.rate, 1)), 0) as total,
			COUNT(*) as count
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		GROUP BY c.id, c.name, c.icon, c.color
		ORDER BY total DESC
		LIMIT ?
	`
	rows, err := h.db.Query(query, baseCurrency, userID, txnType, dateFrom, dateTo, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.CategoryStats
	var grandTotal float64
	for rows.Next() {
		var cs models.CategoryStats
		rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.Icon, &cs.Color, &cs.TotalAmount, &cs.Count)
		categories = append(categories, cs)
		grandTotal += cs.TotalAmount
	}

	for i := range categories {
		categories[i].Percentage = safeDiv(categories[i].TotalAmount, grandTotal) * 100
	}

	writeJSON(w, map[string]interface{}{
		"categories": categories,
		"total":      grandTotal,
		"currency":   baseCurrency,
		"type":       txnType,
		"date_from":  dateFrom,
		"date_to":    dateTo,
	})
}

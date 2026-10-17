package analytics

import (
	"fmt"
	"net/http"
	"strconv"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/models"
)

func (h *Handler) GetTopMerchants(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	limitStr := r.URL.Query().Get("limit")

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "date_from and date_to required", http.StatusBadRequest)
		return
	}
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	baseCurrency := h.cfg.FXBaseCurrency
	if baseCurrency == "" {
		baseCurrency = "JPY"
	}

	query := fmt.Sprintf(`
		SELECT merchant, COALESCE(SUM(t.amount * COALESCE(r.rate, 1)), 0) as total, COUNT(*) as count
		FROM transactions t
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE t.user_id = ? AND t.type = 'expense' AND t.merchant IS NOT NULL AND t.merchant != ''
			AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		GROUP BY merchant
		ORDER BY total DESC
		LIMIT %d
	`, limit)

	rows, err := h.db.Query(query, baseCurrency, userID, dateFrom, dateTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var merchants []models.MerchantStats
	for rows.Next() {
		var ms models.MerchantStats
		rows.Scan(&ms.Merchant, &ms.TotalAmount, &ms.Count)
		merchants = append(merchants, ms)
	}

	writeJSON(w, map[string]interface{}{
		"merchants": merchants,
	})
}

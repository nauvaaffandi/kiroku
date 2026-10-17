package reports

import (
	"time"
)

func (h *Handler) getMonthTotals(monthStart, monthEnd time.Time, accountID, categoryID, baseCurrency string, rates map[string]float64) (float64, float64) {
	args := []interface{}{monthStart.Format(time.RFC3339), monthEnd.Format(time.RFC3339)}
	whereClause := "WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ?"

	if accountID != "" {
		whereClause += " AND (t.from_account_id = ? OR t.to_account_id = ?)"
		args = append(args, accountID, accountID)
	}
	if categoryID != "" {
		whereClause += " AND t.category_id = ?"
		args = append(args, categoryID)
	}

	query := `
		SELECT t.type, t.currency, SUM(t.amount) as total
		FROM transactions t
		` + whereClause + `
		GROUP BY t.type, t.currency
	`

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()

	var income, expense float64
	for rows.Next() {
		var typ, currency string
		var total float64
		rows.Scan(&typ, &currency, &total)
		rate := rates[currency]
		if rate == 0 {
			rate = 1
		}
		converted := total * rate
		if typ == "income" {
			income += converted
		} else if typ == "expense" {
			expense += converted
		}
	}

	return income, expense
}

func (h *Handler) getFxRates(baseCurrency string) map[string]float64 {
	rates := make(map[string]float64)
	rows, err := h.db.Query("SELECT base, rate FROM fx_rates WHERE quote = ?", baseCurrency)
	if err != nil {
		return rates
	}
	defer rows.Close()

	for rows.Next() {
		var from string
		var rate float64
		rows.Scan(&from, &rate)
		rates[from] = rate
	}
	rates[baseCurrency] = 1.0
	return rates
}

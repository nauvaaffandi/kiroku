package analytics

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"kiroku/backend/internal/models"
)

func (h *Handler) getAccountBalances(userID, accountID, baseCurrency string) ([]models.AccountBalance, error) {
	query := `
		SELECT a.id, a.name, a.currency,
			COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount * COALESCE(r.rate, 1) ELSE 0 END), 0) as balance
		FROM accounts a
		LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE a.user_id = ? AND a.archived_at IS NULL
	`
	args := []interface{}{baseCurrency, userID}
	if accountID != "" {
		query += " AND a.id = ?"
		args = append(args, accountID)
	}
	query += " GROUP BY a.id, a.name, a.currency"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.AccountBalance
	for rows.Next() {
		var acc models.AccountBalance
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance); err != nil {
			return nil, err
		}
		if acc.Currency != baseCurrency {
			rate := h.getFxRate(acc.Currency, baseCurrency)
			acc.Balance *= rate
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (h *Handler) getScheduledTransactions(userID string, dateFrom, dateTo time.Time, accountID, baseCurrency string) ([]models.ScheduledItem, error) {
	query := `
		SELECT s.id, s.name, s.amount, s.type, s.account_id, s.currency,
			c.name as category_name, s.rrule
		FROM scheduled_transactions s
		LEFT JOIN categories c ON s.category_id = c.id
		LEFT JOIN accounts a ON s.account_id = a.id
		WHERE s.user_id = ? AND s.enabled = 1 AND a.archived_at IS NULL
	`
	args := []interface{}{userID}
	if accountID != "" {
		query += " AND s.account_id = ?"
		args = append(args, accountID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ScheduledItem
	for rows.Next() {
		var item models.ScheduledItem
		var rrule sql.NullString
		var catName sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Type, &item.AccountID, &item.Currency, &catName, &rrule); err != nil {
			return nil, err
		}
		if catName.Valid {
			item.Category = catName.String
		}
		item.Recurring = rrule.Valid && rrule.String != ""

		dates := h.expandRRULE(rrule.String, dateFrom, dateTo)
		for _, d := range dates {
			item.Date = d.Format("2006-01-02")
			items = append(items, item)
		}
	}

	filtered := make([]models.ScheduledItem, 0)
	for _, item := range items {
		if item.Currency != baseCurrency {
			rate := h.getFxRate(item.Currency, baseCurrency)
			item.Amount *= rate
			item.Currency = baseCurrency
		}
		filtered = append(filtered, item)
	}
	return filtered, nil
}

func (h *Handler) getFutureTransactions(userID string, dateFrom, dateTo time.Time, accountID, baseCurrency string) ([]models.ScheduledItem, error) {
	query := `
		SELECT t.id, t.note, t.amount, t.type, t.from_account_id, t.currency, c.name as category_name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN accounts a ON t.from_account_id = a.id
		WHERE t.user_id = ? AND t.deleted_at IS NULL AND t.date_time > datetime('now')
			AND t.date_time >= ? AND t.date_time <= ?
			AND a.archived_at IS NULL
	`
	args := []interface{}{userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02 23:59:59")}
	if accountID != "" {
		query += " AND t.from_account_id = ?"
		args = append(args, accountID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ScheduledItem
	for rows.Next() {
		var item models.ScheduledItem
		var catName sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Type, &item.AccountID, &item.Currency, &catName); err != nil {
			return nil, err
		}
		if catName.Valid {
			item.Category = catName.String
		}
		item.Recurring = false
		item.Date = ""
		if item.Currency != baseCurrency {
			rate := h.getFxRate(item.Currency, baseCurrency)
			item.Amount *= rate
			item.Currency = baseCurrency
		}
		items = append(items, item)
	}
	return items, nil
}

func (h *Handler) getHistoricalTransactions(userID string, dateFrom, dateTo time.Time, accountID, baseCurrency string) ([]models.ScheduledItem, error) {
	query := `
		SELECT t.id, t.note, t.amount, t.type, t.from_account_id, t.date_time, t.currency, c.name as category_name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN accounts a ON t.from_account_id = a.id
		WHERE t.user_id = ? AND t.deleted_at IS NULL AND t.date_time < datetime('now')
			AND t.date_time >= ? AND t.date_time <= ?
			AND a.archived_at IS NULL
	`
	args := []interface{}{userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02 23:59:59")}
	if accountID != "" {
		query += " AND t.from_account_id = ?"
		args = append(args, accountID)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ScheduledItem
	for rows.Next() {
		var item models.ScheduledItem
		var catName sql.NullString
		var dateTimeStr string
		if err := rows.Scan(&item.ID, &item.Name, &item.Amount, &item.Type, &item.AccountID, &dateTimeStr, &item.Currency, &catName); err != nil {
			return nil, err
		}
		if catName.Valid {
			item.Category = catName.String
		}
		item.Recurring = false
		if dt, err := time.Parse("2006-01-02 15:04:05", dateTimeStr); err == nil {
			item.Date = dt.Format("2006-01-02")
		}
		if item.Currency != baseCurrency {
			rate := h.getFxRate(item.Currency, baseCurrency)
			item.Amount *= rate
			item.Currency = baseCurrency
		}
		items = append(items, item)
	}
	return items, nil
}

func (h *Handler) calculateBalanceForRange(userID, dateFrom, dateTo, accountID string) (float64, float64, float64) {
	query := `SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END), 0) FROM transactions WHERE user_id = ? AND date >= ? AND date < ? AND deleted_at IS NULL`
	args := []interface{}{userID, dateFrom, dateTo}
	if accountID != "" {
		query += ` AND from_account_id = ?`
		args = append(args, accountID)
	}
	var net float64
	h.db.QueryRow(query, args...).Scan(&net)
	return 0, net, net
}

func (h *Handler) calculateBalanceForAccountRange(userID string, accountID int64, dateFrom, dateTo string) (float64, float64, float64) {
	var opening, net float64
	h.db.QueryRow(`SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END), 0) FROM transactions WHERE user_id = ? AND from_account_id = ? AND date < ? AND deleted_at IS NULL`, userID, accountID, dateFrom).Scan(&opening)
	h.db.QueryRow(`SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END), 0) FROM transactions WHERE user_id = ? AND from_account_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL`, userID, accountID, dateFrom, dateTo).Scan(&net)
	return opening, opening + net, net
}

func (h *Handler) calculateBalanceForCurrencyRange(userID, currency, dateFrom, dateTo string) (float64, float64, float64) {
	var opening, net float64
	h.db.QueryRow(`SELECT COALESCE(SUM(t.amount * CASE WHEN t.type = 'income' THEN 1 WHEN t.type = 'expense' THEN -1 ELSE 0 END), 0) FROM transactions t JOIN accounts a ON t.from_account_id = a.id WHERE a.user_id = ? AND a.currency = ? AND t.date < ? AND t.deleted_at IS NULL`, userID, currency, dateFrom).Scan(&opening)
	h.db.QueryRow(`SELECT COALESCE(SUM(t.amount * CASE WHEN t.type = 'income' THEN 1 WHEN t.type = 'expense' THEN -1 ELSE 0 END), 0) FROM transactions t JOIN accounts a ON t.from_account_id = a.id WHERE a.user_id = ? AND a.currency = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL`, userID, currency, dateFrom, dateTo).Scan(&net)
	return opening, opening + net, net
}

func (h *Handler) getPeriodStats(userID, dateFrom, dateTo string) map[string]float64 {
	var income, expense float64
	var count int
	h.db.QueryRow(`SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0), COUNT(*) FROM transactions WHERE user_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL`, userID, dateFrom, dateTo).Scan(&income, &expense, &count)
	return map[string]float64{
		"income":  income,
		"expense": expense,
		"net":     income - expense,
		"count":   float64(count),
	}
}

func calculateDays(dateFrom, dateTo string) float64 {
	start, _ := time.Parse("2006-01-02", dateFrom)
	end, _ := time.Parse("2006-01-02", dateTo)
	return end.Sub(start).Hours()/24 + 1
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

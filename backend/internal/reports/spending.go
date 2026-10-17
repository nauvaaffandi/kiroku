package reports

import (
	"database/sql"

	"kiroku/backend/internal/models"
)

func (h *Handler) fetchChildren(parent *models.Group, dateFrom, dateTo, accountID, baseCurrency string) {
	args := []interface{}{parent.ID, dateFrom, dateTo}
	if accountID != "" {
		args = append(args, accountID, accountID)
	}

	query := `
		SELECT 
			c.id, c.name, c.color, c.kind,
			SUM(t.amount * COALESCE(r.rate, 1)) as total_amount,
			COUNT(*) as count,
			? as currency
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN fx_rates r ON t.currency = r.base AND r.quote = ?
		WHERE t.deleted_at IS NULL AND t.date_time >= ? AND t.date_time <= ? AND c.parent_id = ?
	`
	if accountID != "" {
		query += " AND (t.from_account_id = ? OR t.to_account_id = ?)"
	}
	query += " GROUP BY c.id, c.name, c.color, c.kind ORDER BY total_amount DESC"

	args = append(args, baseCurrency)
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Group
		var kind sql.NullString
		var color sql.NullString
		if err := rows.Scan(&g.ID, &g.Name, &color, &kind, &g.Amount, &g.Count, &g.Currency); err != nil {
			continue
		}
		if color.Valid {
			g.Color = color.String
		}
		if kind.Valid {
			g.Kind = kind.String
		}
		parent.Children = append(parent.Children, g)
	}
}

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (r *ImpexpRepo) ExportCSV(ctx context.Context, dateFrom, dateTo, typ, categoryID, accountID string) ([][]string, error) {
	query := `SELECT t.date_time, t.type, t.amount, t.currency,
		fa.name as from_account, ta.name as to_account,
		c.name as category, t.labels, t.note, t.payee_payer,
		t.payment_type, t.warranty_until, t.status,
		t.place_name, t.attachments
		FROM transactions t
		LEFT JOIN accounts fa ON t.from_account_id = fa.id
		LEFT JOIN accounts ta ON t.to_account_id = ta.id
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.deleted_at IS NULL`
	args := []interface{}{}
	conditions := []string{}
	if dateFrom != "" {
		conditions = append(conditions, "t.date_time >= ?")
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		conditions = append(conditions, "t.date_time <= ?")
		args = append(args, dateTo)
	}
	if typ != "" {
		conditions = append(conditions, "t.type = ?")
		args = append(args, typ)
	}
	if categoryID != "" {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, categoryID)
	}
	if accountID != "" {
		conditions = append(conditions, "(t.from_account_id = ? OR t.to_account_id = ?)")
		args = append(args, accountID, accountID)
	}
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.date_time DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	headers := []string{"date_time", "type", "amount", "currency", "from_account", "to_account", "category", "labels", "note", "payee_payer", "payment_type", "warranty_until", "status", "place", "attachments"}
	records := [][]string{headers}

	for rows.Next() {
		var dateTime time.Time
		var typStr, currency, fromAcc, toAcc, category, labels, note, payeePayer, paymentType, warrantyUntil, status, placeName, attachments sql.NullString
		var amount int64
		if err := rows.Scan(&dateTime, &typStr, &amount, &currency, &fromAcc, &toAcc, &category, &labels, &note, &payeePayer, &paymentType, &warrantyUntil, &status, &placeName, &attachments); err != nil {
			continue
		}
		record := []string{
			dateTime.Format(time.RFC3339),
			typStr.String,
			fmt.Sprintf("%d", amount),
			currency.String,
			fromAcc.String,
			toAcc.String,
			category.String,
			labels.String,
			note.String,
			payeePayer.String,
			paymentType.String,
			warrantyUntil.String,
			status.String,
			placeName.String,
			attachments.String,
		}
		records = append(records, record)
	}
	return records, rows.Err()
}
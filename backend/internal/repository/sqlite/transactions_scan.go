package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"kiroku/backend/internal/domain"
)

func (r *TransactionsRepo) scanTransactions(ctx context.Context, query string, args ...interface{}) ([]domain.Transaction, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		var fromAcc, toAcc, catID, templateID sql.NullString
		var warranty sql.NullTime
		var placeLat, placeLng sql.NullFloat64
		var labelsStr, attachmentsStr sql.NullString
		if err := rows.Scan(&t.ID, &t.Type, &t.DateTime, &t.Amount, &t.Currency, &fromAcc, &toAcc, &catID, &labelsStr, &t.Note, &t.PayeePayer, &t.PaymentType, &warranty, &t.Status, &t.PlaceName, &placeLat, &placeLng, &attachmentsStr, &templateID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if fromAcc.Valid {
			t.FromAccountID = &fromAcc.String
		}
		if toAcc.Valid {
			t.ToAccountID = &toAcc.String
		}
		if catID.Valid {
			t.CategoryID = &catID.String
		}
		if templateID.Valid {
			t.TemplateID = &templateID.String
		}
		if warranty.Valid {
			t.WarrantyUntil = &warranty.Time
		}
		if placeLat.Valid {
			t.PlaceLat = &placeLat.Float64
		}
		if placeLng.Valid {
			t.PlaceLng = &placeLng.Float64
		}
		if labelsStr.Valid && labelsStr.String != "" {
			t.Labels = strings.Split(labelsStr.String, ",")
		}
		// Load label IDs from junction table
		labelRows, err := r.db.QueryContext(ctx, `SELECT label_id FROM transaction_labels WHERE transaction_id = ?`, t.ID)
		if err == nil {
			var lIDs []string
			for labelRows.Next() {
				var lid string
				if err := labelRows.Scan(&lid); err == nil {
					lIDs = append(lIDs, lid)
				}
			}
			labelRows.Close()
			if lIDs != nil {
				t.LabelIDs = lIDs
			} else {
				t.LabelIDs = []string{}
			}
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
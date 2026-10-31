package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"kiroku/backend/internal/domain"
)

func (r *ImpexpRepo) ImportPreview(ctx context.Context, headers []string, dataRows [][]string) (*domain.ImportPreview, error) {
	if len(dataRows) > 10 {
		dataRows = dataRows[:10]
	}

	preview := &domain.ImportPreview{PreviewID: genID()}
	catCache := map[string]string{}
	accCache := map[string]string{}

	for _, row := range dataRows {
		if len(row) < len(headers) {
			continue
		}
		rec := rowToImport(headers, row)
		preview.Rows = append(preview.Rows, rec)

		if rec.Category != "" {
			if _, ok := catCache[rec.Category]; !ok {
				match := r.fuzzyMatchCategory(rec.Category)
				preview.Categories = append(preview.Categories, match)
				catCache[rec.Category] = match.MatchID
			}
		}
		if rec.FromAccount != "" {
			if _, ok := accCache[rec.FromAccount]; !ok {
				match := r.fuzzyMatchAccount(rec.FromAccount)
				preview.Accounts = append(preview.Accounts, match)
				accCache[rec.FromAccount] = match.MatchID
			}
		}
		if rec.ToAccount != "" {
			if _, ok := accCache[rec.ToAccount]; !ok {
				match := r.fuzzyMatchAccount(rec.ToAccount)
				preview.Accounts = append(preview.Accounts, match)
				accCache[rec.ToAccount] = match.MatchID
			}
		}
	}

	matches := struct {
		Categories []domain.CategoryMatch `json:"categories"`
		Accounts   []domain.AccountMatch  `json:"accounts"`
	}{
		Categories: preview.Categories,
		Accounts:   preview.Accounts,
	}
	r.storePreview(preview.PreviewID, preview.Rows, matches)
	return preview, nil
}

func (r *ImpexpRepo) GetPreview(ctx context.Context, previewID string) (*domain.ImportPreview, error) {
	var data string
	err := r.db.QueryRowContext(ctx, `SELECT data FROM import_log WHERE id = ?`, previewID).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var preview domain.ImportPreview
	if err := json.Unmarshal([]byte(data), &preview); err != nil {
		return nil, err
	}
	return &preview, nil
}

func (r *ImpexpRepo) Confirm(ctx context.Context, previewID string) (int, int, []string, error) {
	preview, err := r.GetPreview(ctx, previewID)
	if err != nil {
		return 0, 0, nil, err
	}
	if preview == nil {
		return 0, 0, nil, sql.ErrNoRows
	}

	created := 0
	skipped := 0
	var errors []string

	for _, row := range preview.Rows {
		hash := rowHash(row)
		var exists bool
		r.db.QueryRowContext(ctx, `SELECT 1 FROM import_log WHERE id = ?`, hash).Scan(&exists)
		if exists {
			skipped++
			continue
		}

		dateTime, err := time.Parse(time.RFC3339, row.DateTime)
		if err != nil {
			dateTime = time.Now()
		}

		if row.Amount == 0 {
			errors = append(errors, row.DateTime+": amount required")
			continue
		}

		fromAccID := resolveAccountByName(r.db, row.FromAccount)
		toAccID := resolveAccountByName(r.db, row.ToAccount)
		catID := resolveCategoryByName(r.db, row.Category)

		if row.Type == "expense" && fromAccID == nil {
			errors = append(errors, row.DateTime+": expense requires from_account")
			continue
		}
		if row.Type == "income" && toAccID == nil {
			errors = append(errors, row.DateTime+": income requires to_account")
			continue
		}

		id := genID()
		warranty := parseTime(row.WarrantyUntil)
		status := defaultStatus(row.Status)
		_, err = r.db.ExecContext(ctx, `
			INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, attachments)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, row.Type, dateTime, row.Amount, row.Currency, fromAccID, toAccID, catID, row.Labels, row.Note, row.PayeePayer, row.PaymentType, warranty, status, row.Place, row.Attachments)
		if err != nil {
			errors = append(errors, row.DateTime+": "+err.Error())
			continue
		}

		r.db.ExecContext(ctx, `INSERT INTO import_log (id, data, created_at) VALUES (?, ?, ?)`, hash, "{}", time.Now())
		created++
	}

	return created, skipped, errors, nil
}
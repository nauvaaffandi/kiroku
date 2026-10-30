package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type TransactionsRepo struct{ db *db.DB }

func NewTransactionsRepo(database *db.DB) *TransactionsRepo { return &TransactionsRepo{db: database} }

func (r *TransactionsRepo) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, error) {
	query := `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE deleted_at IS NULL`
	args := []interface{}{}
	conditions := []string{}

	if filter.AccountID != nil && *filter.AccountID != "" {
		conditions = append(conditions, "(from_account_id = ? OR to_account_id = ?)")
		args = append(args, *filter.AccountID, *filter.AccountID)
	}
	if filter.CategoryID != nil && *filter.CategoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, *filter.CategoryID)
	}
	if len(filter.Types) > 0 {
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Types[0])
	}
	if filter.Status != nil && *filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, *filter.Status)
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, "date_time >= ?")
		args = append(args, filter.DateFrom.Format(time.RFC3339))
	}
	if filter.DateTo != nil {
		conditions = append(conditions, "date_time <= ?")
		args = append(args, filter.DateTo.Format(time.RFC3339))
	}
	if filter.MinAmount != nil {
		conditions = append(conditions, "amount >= ?")
		args = append(args, *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		conditions = append(conditions, "amount <= ?")
		args = append(args, *filter.MaxAmount)
	}
	if filter.SearchQuery != "" {
		conditions = append(conditions, "id IN (SELECT id FROM transactions_fts WHERE transactions_fts MATCH ?)")
		args = append(args, filter.SearchQuery)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	sortBy := filter.Sort
	if sortBy == "" {
		sortBy = "date_time"
	}
	order := filter.Order
	if order == "" {
		order = "DESC"
	}
	query += " ORDER BY " + sortBy + " " + order

	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	return r.scanTransactions(ctx, query, args...)
}

func (r *TransactionsRepo) Count(ctx context.Context, filter domain.TransactionFilter) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *TransactionsRepo) Get(ctx context.Context, id string) (*domain.Transaction, error) {
	list, err := r.scanTransactions(ctx, `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func (r *TransactionsRepo) Create(ctx context.Context, input domain.CreateTransactionInput) (*domain.Transaction, error) {
	id := genID()
	status := input.Status
	if status == "" {
		status = "cleared"
	}
	labelsStr := strings.Join(input.Labels, ",")
	attachmentsStr, _ := jsonMarshal(input.Attachments)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.Type, input.DateTime, input.Amount, input.Currency, input.FromAccountID, input.ToAccountID, input.CategoryID, labelsStr, input.Note, input.PayeePayer, input.PaymentType, input.WarrantyUntil, status, input.PlaceName, input.PlaceLat, input.PlaceLng, string(attachmentsStr), input.TemplateID)
	if err != nil {
		return nil, err
	}

	if len(input.LabelIDs) > 0 {
		for _, lid := range input.LabelIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO transaction_labels (transaction_id, label_id) VALUES (?, ?)`, id, lid); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *TransactionsRepo) Update(ctx context.Context, id string, input domain.UpdateTransactionInput) (*domain.Transaction, error) {
	setParts := []string{}
	args := []interface{}{}
	if input.DateTime != nil {
		setParts = append(setParts, "date_time = ?")
		args = append(args, input.DateTime.Format(time.RFC3339))
	}
	if input.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *input.Amount)
	}
	if input.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *input.CategoryID)
	}
	if input.Labels != nil {
		setParts = append(setParts, "labels = ?")
		args = append(args, strings.Join(input.Labels, ","))
	}
	if input.Note != nil {
		setParts = append(setParts, "note = ?")
		args = append(args, *input.Note)
	}
	if input.PayeePayer != nil {
		setParts = append(setParts, "payee_payer = ?")
		args = append(args, *input.PayeePayer)
	}
	if input.PaymentType != nil {
		setParts = append(setParts, "payment_type = ?")
		args = append(args, *input.PaymentType)
	}
	if input.WarrantyUntil != nil {
		setParts = append(setParts, "warranty_until = ?")
		args = append(args, input.WarrantyUntil.Format("2006-01-02"))
	}
	if input.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *input.Status)
	}
	if input.PlaceName != nil {
		setParts = append(setParts, "place_name = ?")
		args = append(args, *input.PlaceName)
	}
	if input.PlaceLat != nil {
		setParts = append(setParts, "place_lat = ?")
		args = append(args, *input.PlaceLat)
	}
	if input.PlaceLng != nil {
		setParts = append(setParts, "place_lng = ?")
		args = append(args, *input.PlaceLng)
	}
	if input.Attachments != nil {
		attStr, _ := jsonMarshal(input.Attachments)
		setParts = append(setParts, "attachments = ?")
		args = append(args, string(attStr))
	}
	if len(setParts) == 0 && input.LabelIDs == nil {
		return r.Get(ctx, id)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if len(setParts) > 0 {
		args = append(args, id)
		_, err := tx.ExecContext(ctx, "UPDATE transactions SET "+strings.Join(setParts, ", ")+" WHERE id = ?", args...)
		if err != nil {
			return nil, err
		}
	}

	if input.LabelIDs != nil {
		_, err := tx.ExecContext(ctx, `DELETE FROM transaction_labels WHERE transaction_id = ?`, id)
		if err != nil {
			return nil, err
		}
		for _, lid := range *input.LabelIDs {
			_, err := tx.ExecContext(ctx, `INSERT INTO transaction_labels (transaction_id, label_id) VALUES (?, ?)`, id, lid)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *TransactionsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *TransactionsRepo) Search(ctx context.Context, query string, limit int) ([]domain.Transaction, error) {
	if limit <= 0 {
		limit = 50
	}
	q := `SELECT id, type, date_time, amount, currency, from_account_id, to_account_id, category_id, labels, note, payee_payer, payment_type, warranty_until, status, place_name, place_lat, place_lng, attachments, template_id, created_at, updated_at FROM transactions WHERE id IN (SELECT id FROM transactions_fts WHERE transactions_fts MATCH ?) AND deleted_at IS NULL ORDER BY date_time DESC LIMIT ?`
	return r.scanTransactions(ctx, q, query, limit)
}

func (r *TransactionsRepo) Autocomplete(ctx context.Context, query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	q := `SELECT DISTINCT note FROM transactions WHERE note LIKE ? AND deleted_at IS NULL ORDER BY date_time DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, q, query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var note string
		if err := rows.Scan(&note); err == nil && note != "" {
			out = append(out, note)
		}
	}
	return out, rows.Err()
}

func (r *TransactionsRepo) Summary(ctx context.Context, filter domain.TransactionFilter) ([]domain.SummaryItem, error) {
	query := `SELECT currency, SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income, SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense, SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END) as net FROM transactions WHERE deleted_at IS NULL`
	args := []interface{}{}
	conditions := []string{}

	if filter.CategoryID != nil && *filter.CategoryID != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, *filter.CategoryID)
	}
	if filter.AccountID != nil && *filter.AccountID != "" {
		conditions = append(conditions, "(from_account_id = ? OR to_account_id = ?)")
		args = append(args, *filter.AccountID, *filter.AccountID)
	}
	if len(filter.Types) > 0 {
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Types[0])
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, "date_time >= ?")
		args = append(args, filter.DateFrom.Format(time.RFC3339))
	}
	if filter.DateTo != nil {
		conditions = append(conditions, "date_time <= ?")
		args = append(args, filter.DateTo.Format(time.RFC3339))
	}
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += " GROUP BY currency"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summary []domain.SummaryItem
	for rows.Next() {
		var s domain.SummaryItem
		if err := rows.Scan(&s.Currency, &s.Income, &s.Expense, &s.Net); err != nil {
			return nil, err
		}
		summary = append(summary, s)
	}
	return summary, rows.Err()
}

func (r *TransactionsRepo) Batch(ctx context.Context, input domain.BatchInput) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	affected := 0
	for _, id := range input.IDs {
		var res sql.Result
		var execErr error
		switch input.Action {
		case domain.BatchDelete:
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", id)
		case domain.BatchRecategorize:
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET category_id = ? WHERE id = ? AND deleted_at IS NULL", input.CategoryID, id)
		case domain.BatchRelabel:
			labels := strings.Join(input.Labels, ",")
			res, execErr = tx.ExecContext(ctx, "UPDATE transactions SET labels = ? WHERE id = ? AND deleted_at IS NULL", labels, id)
		default:
			tx.Rollback()
			return 0, fmt.Errorf("unknown action: %s", input.Action)
		}
		if execErr != nil {
			continue
		}
		if n, _ := res.RowsAffected(); n > 0 {
			affected++
		}
	}
	return affected, tx.Commit()
}

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

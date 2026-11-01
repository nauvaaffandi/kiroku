package sqlite

import (
	"context"
	"strings"
	"time"

	"kiroku/backend/internal/domain"
)

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

	query += " ORDER BY " + sanitizeSortOrder(filter.Sort, filter.Order)

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

// allowedSortColumns lists every column selected by List that may be used
// for ordering. User-supplied values not present here fall back to date_time.
var allowedSortColumns = map[string]struct{}{
	"id":             {},
	"type":           {},
	"date_time":      {},
	"amount":         {},
	"currency":       {},
	"from_account_id": {},
	"to_account_id":  {},
	"category_id":    {},
	"labels":         {},
	"note":           {},
	"payee_payer":    {},
	"payment_type":   {},
	"warranty_until": {},
	"status":         {},
	"place_name":     {},
	"place_lat":      {},
	"place_lng":      {},
	"attachments":    {},
	"template_id":    {},
	"created_at":     {},
	"updated_at":     {},
}

// allowedDirections lists safe ORDER BY directions (case-insensitive).
var allowedDirections = map[string]struct{}{
	"ASC":  {},
	"DESC": {},
}

// sanitizeSortOrder returns a safe ORDER BY fragment whitelisted against
// SQL injection from user-supplied sort/order query params.
func sanitizeSortOrder(sortBy, order string) string {
	if _, ok := allowedSortColumns[sortBy]; !ok {
		sortBy = "date_time"
	}
	o := strings.ToUpper(order)
	if _, ok := allowedDirections[o]; !ok {
		o = "DESC"
	}
	return sortBy + " " + o
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
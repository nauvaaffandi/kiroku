package sqlite

import (
	"context"
	"database/sql"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type SharedRepo struct{ db *db.DB }

func NewSharedRepo(database *db.DB) *SharedRepo { return &SharedRepo{db: database} }

func (r *SharedRepo) Create(ctx context.Context, userID, accountID string, expiresAt *time.Time) (*domain.Share, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounts WHERE id = ? AND user_id = ?`, accountID, userID).Scan(&count)
	if err != nil || count == 0 {
		return nil, nil
	}

	var expiresAtNull sql.NullTime
	if expiresAt != nil {
		expiresAtNull = sql.NullTime{Time: *expiresAt, Valid: true}
	}

	shareID := genID()
	token := genID()

	_, err = r.db.ExecContext(ctx, `INSERT INTO shares (id, account_id, token, expires_at) VALUES (?, ?, ?, ?)`, shareID, accountID, token, expiresAtNull)
	if err != nil {
		return nil, err
	}

	var accountName string
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM accounts WHERE id = ?`, accountID).Scan(&accountName)

	return &domain.Share{
		ID:          shareID,
		AccountID:   accountID,
		AccountName: accountName,
		Token:       token,
		ExpiresAt:   expiresAt,
	}, nil
}

func (r *SharedRepo) List(ctx context.Context, userID string) ([]domain.Share, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.account_id, a.name, s.token, s.expires_at, s.created_at, s.last_accessed_at
		FROM shares s JOIN accounts a ON s.account_id = a.id
		WHERE a.user_id = ? ORDER BY s.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Share
	for rows.Next() {
		var s domain.Share
		var token string
		var expiresAt, lastAccessedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.AccountID, &s.AccountName, &token, &expiresAt, &s.CreatedAt, &lastAccessedAt); err != nil {
			continue
		}
		if len(token) > 8 {
			s.Token = token[:8] + "..."
		} else {
			s.Token = token
		}
		if expiresAt.Valid {
			s.ExpiresAt = &expiresAt.Time
		}
		if lastAccessedAt.Valid {
			s.LastAccessedAt = &lastAccessedAt.Time
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SharedRepo) Revoke(ctx context.Context, userID, shareID string) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM shares WHERE id = ? AND account_id IN (SELECT id FROM accounts WHERE user_id = ?)`, shareID, userID)
	if err != nil {
		return false, err
	}
	affected, _ := result.RowsAffected()
	return affected > 0, nil
}

func (r *SharedRepo) GetByToken(ctx context.Context, token string) (*domain.Share, error) {
	var s domain.Share
	var expiresAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT id, account_id, expires_at FROM shares WHERE token = ?`, token).Scan(&s.ID, &s.AccountID, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		s.ExpiresAt = &expiresAt.Time
		if s.ExpiresAt.Before(time.Now()) {
			return nil, sql.ErrNoRows
		}
	}
	_, _ = r.db.ExecContext(ctx, "UPDATE shares SET last_accessed_at = CURRENT_TIMESTAMP WHERE id = ?", s.ID)
	return &s, nil
}

func (r *SharedRepo) TouchAccessed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE shares SET last_accessed_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *SharedRepo) GetAccountInfo(ctx context.Context, accountID string) (*domain.AccountWithBalance, error) {
	var balance float64
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, currency,
			COALESCE((SELECT SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END) FROM transactions WHERE account_id = a.id), 0) as balance
		FROM accounts a WHERE id = ?`, accountID).Scan(&balance)
	if err != nil {
		return nil, err
	}
	return &domain.AccountWithBalance{Account: domain.Account{ID: accountID}, Balance: balance}, nil
}

func (r *SharedRepo) GetTransactions(ctx context.Context, accountID string) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.date, t.amount, t.type, t.note, t.payee,
			c.name, c.icon, c.color
		FROM transactions t LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.account_id = ?
		ORDER BY t.date DESC, t.created_at DESC LIMIT 100`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []map[string]interface{}
	for rows.Next() {
		var tx struct {
			ID            string
			Date          string
			Amount        float64
			Type          string
			Note          string
			Payee         string
			Category      string
			CategoryIcon  sql.NullString
			CategoryColor sql.NullString
		}
		if err := rows.Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Type, &tx.Note, &tx.Payee, &tx.Category, &tx.CategoryIcon, &tx.CategoryColor); err != nil {
			continue
		}
		entry := map[string]interface{}{
			"id":            tx.ID,
			"date":          tx.Date,
			"amount":        tx.Amount,
			"type":          tx.Type,
			"note":          tx.Note,
			"payee":         tx.Payee,
			"category":      tx.Category,
			"category_icon": tx.CategoryIcon.String,
			"category_color": tx.CategoryColor.String,
		}
		out = append(out, entry)
	}
	return out, rows.Err()
}

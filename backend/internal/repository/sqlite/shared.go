package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type SharedRepo struct{ db *db.DB }

func NewSharedRepo(database *db.DB) *SharedRepo { return &SharedRepo{db: database} }

func (r *SharedRepo) Create(ctx context.Context, userID, accountID string, expiresAt *time.Time) (*domain.Share, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounts WHERE id = ?`, accountID).Scan(&count)
	if err != nil {
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "Create").
			Str("account_id", accountID).
			Msg("account lookup failed")
		return nil, err
	}
	if count == 0 {
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
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "Create").
			Str("share_id", shareID).
			Str("account_id", accountID).
			Msg("insert share failed")
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
		ORDER BY s.created_at DESC`)
	if err != nil {
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "List").
			Msg("list shares query failed")
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Share, 0)
	for rows.Next() {
		var s domain.Share
		var token string
		var expiresAt, lastAccessedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.AccountID, &s.AccountName, &token, &expiresAt, &s.CreatedAt, &lastAccessedAt); err != nil {
			log.Error().
				Err(err).
				Str("component", "shared_repo").
				Str("method", "List").
				Msg("share row scan failed")
			continue
		}
		s.Token = token
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
	result, err := r.db.ExecContext(ctx, `DELETE FROM shares WHERE id = ?`, shareID)
	if err != nil {
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "Revoke").
			Str("share_id", shareID).
			Msg("revoke share query failed")
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
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "GetByToken").
			Str("token", token).
			Msg("get share by token failed")
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
	query := `SELECT a.id, a.name, a.currency, a.type, a.icon, a.color, a.archived,
		COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
		COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a LEFT JOIN transactions t ON (t.from_account_id = a.id OR t.to_account_id = a.id) AND t.deleted_at IS NULL
		WHERE a.id = ? AND a.archived = FALSE GROUP BY a.id`
	var a domain.AccountWithBalance
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&a.ID, &a.Name, &a.Currency, &a.Type, &a.Icon, &a.Color, &a.Archived, &a.Balance)
	if err != nil {
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "GetAccountInfo").
			Str("account_id", accountID).
			Msg("get account info failed")
		return nil, err
	}
	return &a, nil
}

func (r *SharedRepo) GetTransactions(ctx context.Context, accountID string) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.date_time, t.amount, t.type, t.note, t.payee_payer,
			COALESCE(c.name, ''), COALESCE(c.icon, ''), COALESCE(c.color, '')
		FROM transactions t LEFT JOIN categories c ON t.category_id = c.id
		WHERE (t.from_account_id = ? OR t.to_account_id = ?) AND t.deleted_at IS NULL
		ORDER BY t.date_time DESC, t.created_at DESC LIMIT 100`, accountID, accountID)
	if err != nil {
		log.Error().
			Err(err).
			Str("component", "shared_repo").
			Str("method", "GetTransactions").
			Str("account_id", accountID).
			Msg("get transactions failed")
		return nil, err
	}
	defer rows.Close()

	out := make([]map[string]interface{}, 0)
	for rows.Next() {
		var tx struct {
			ID            string
			DateTime      string
			Amount        int64
			Type          string
			Note          string
			PayeePayer    string
			Category      string
			CategoryIcon  string
			CategoryColor string
		}
		if err := rows.Scan(&tx.ID, &tx.DateTime, &tx.Amount, &tx.Type, &tx.Note, &tx.PayeePayer, &tx.Category, &tx.CategoryIcon, &tx.CategoryColor); err != nil {
			log.Error().
				Err(err).
				Str("component", "shared_repo").
				Str("method", "GetTransactions").
				Msg("scan transaction failed")
			continue
		}
		entry := map[string]interface{}{
			"id":             tx.ID,
			"date":           tx.DateTime,
			"amount":         float64(tx.Amount) / 100.0,
			"type":           tx.Type,
			"note":           tx.Note,
			"payee":          tx.PayeePayer,
			"category":       tx.Category,
			"category_icon":  tx.CategoryIcon,
			"category_color": tx.CategoryColor,
		}
		out = append(out, entry)
	}
	return out, rows.Err()
}

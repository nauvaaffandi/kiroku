package sqlite

import (
	"context"
	"database/sql"
	"time"

	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/db"
)

type UserRepo struct{ db *db.DB }

func NewUserRepo(database *db.DB) *UserRepo { return &UserRepo{db: database} }

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, pin_hash, created_at, updated_at FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Email, &u.PinHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, pin_hash, created_at, updated_at FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Email, &u.PinHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdatePinHash(ctx context.Context, id, pinHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET pin_hash = ? WHERE id = ?`, pinHash, id)
	return err
}

type PinChangeTokenRepo struct{ db *db.DB }

func NewPinChangeTokenRepo(database *db.DB) *PinChangeTokenRepo { return &PinChangeTokenRepo{db: database} }

func (r *PinChangeTokenRepo) Create(ctx context.Context, token, userID, purpose string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at) VALUES (?, ?, ?, ?, ?)`, genID(), userID, token, purpose, expiresAt)
	return err
}

func (r *PinChangeTokenRepo) Get(ctx context.Context, token string) (userID string, purpose string, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT user_id, type FROM pin_change_tokens WHERE token = ? AND expires_at > datetime('now')`, token).Scan(&userID, &purpose)
	return
}

func (r *PinChangeTokenRepo) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pin_change_tokens WHERE token = ?`, token)
	return err
}

func (r *PinChangeTokenRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pin_change_tokens WHERE expires_at <= datetime('now')`)
	return err
}

type TelegramBindingRepo struct{ db *db.DB }

func NewTelegramBindingRepo(database *db.DB) *TelegramBindingRepo { return &TelegramBindingRepo{db: database} }

func (r *TelegramBindingRepo) Get(ctx context.Context, chatID int64) (*domain.TelegramBinding, error) {
	var b domain.TelegramBinding
	err := r.db.QueryRowContext(ctx, `SELECT chat_id, user_id, username FROM telegram_bindings WHERE telegram_chat_id = ?`, chatID).Scan(&b.ChatID, &b.UserID, &b.Username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *TelegramBindingRepo) Upsert(ctx context.Context, binding domain.TelegramBinding) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO telegram_bindings (id, user_id, telegram_chat_id, telegram_user_id) VALUES (?, ?, ?, ?) ON CONFLICT(telegram_chat_id) DO UPDATE SET user_id = excluded.user_id, telegram_user_id = excluded.telegram_user_id`, genID(), binding.UserID, binding.ChatID, binding.Username)
	return err
}

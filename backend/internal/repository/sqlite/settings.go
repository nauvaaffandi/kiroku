package sqlite

import (
	"context"

	"kiroku/backend/internal/db"
)

type SettingsRepo struct{ db *db.DB }

func NewSettingsRepo(database *db.DB) *SettingsRepo { return &SettingsRepo{db: database} }

func (r *SettingsRepo) GetAll(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, rows.Err()
}

func (r *SettingsRepo) GetOne(ctx context.Context, key string) (string, error) {
	var v string
	err := r.db.QueryRowContext(ctx, "SELECT value FROM settings WHERE key = ?", key).Scan(&v)
	return v, err
}

func (r *SettingsRepo) SetOne(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	return err
}

func (r *SettingsRepo) UpdateAll(ctx context.Context, settings map[string]string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for k, v := range settings {
		if _, err := tx.ExecContext(ctx, `INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, k, v); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

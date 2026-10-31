package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type PushSubscriptionRepo struct{ db *db.DB }

func NewPushSubscriptionRepo(database *db.DB) *PushSubscriptionRepo {
	return &PushSubscriptionRepo{db: database}
}

func (r *PushSubscriptionRepo) Create(ctx context.Context, input domain.CreatePushSubscriptionInput) (*domain.PushSubscription, error) {
	id := fmt.Sprintf("ps_%s", time.Now().Format("20060102150405.000000000"))
	now := time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO push_subscriptions (id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.UserID, input.Endpoint, input.Keys.P256DH, input.Keys.Auth, input.UserAgent, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create push subscription: %w", err)
	}

	return r.Get(ctx, id)
}

func (r *PushSubscriptionRepo) Get(ctx context.Context, id string) (*domain.PushSubscription, error) {
	var s domain.PushSubscription
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at
		FROM push_subscriptions WHERE id = ?
	`, id).Scan(
		&s.ID, &s.UserID, &s.Endpoint, &s.P256DH, &s.Auth, &s.UserAgent, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get push subscription: %w", err)
	}
	return &s, nil
}

func (r *PushSubscriptionRepo) GetByEndpoint(ctx context.Context, endpoint, userID string) (*domain.PushSubscription, error) {
	var s domain.PushSubscription
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at
		FROM push_subscriptions WHERE endpoint = ? AND user_id = ?
	`, endpoint, userID).Scan(
		&s.ID, &s.UserID, &s.Endpoint, &s.P256DH, &s.Auth, &s.UserAgent, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get push subscription by endpoint: %w", err)
	}
	return &s, nil
}

func (r *PushSubscriptionRepo) ListByUser(ctx context.Context, userID string) ([]*domain.PushSubscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at
		FROM push_subscriptions WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list push subscriptions: %w", err)
	}
	defer rows.Close()

	var list []*domain.PushSubscription
	for rows.Next() {
		var s domain.PushSubscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Endpoint, &s.P256DH, &s.Auth, &s.UserAgent, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan push subscription: %w", err)
		}
		list = append(list, &s)
	}
	if list == nil {
		list = []*domain.PushSubscription{}
	}
	return list, rows.Err()
}

func (r *PushSubscriptionRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete push subscription: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PushSubscriptionRepo) DeleteByEndpoint(ctx context.Context, endpoint, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE endpoint = ? AND user_id = ?`, endpoint, userID)
	if err != nil {
		return fmt.Errorf("failed to delete push subscription by endpoint: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) DeleteByUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete push subscriptions by user: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) Update(ctx context.Context, id string, input domain.UpdatePushSubscriptionInput) (*domain.PushSubscription, error) {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE push_subscriptions SET user_agent = COALESCE(?, user_agent), updated_at = ?
		WHERE id = ?
	`, input.UserAgent, now, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update push subscription: %w", err)
	}
	return r.Get(ctx, id)
}

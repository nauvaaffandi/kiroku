package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

const DeliveryChannelWebPush = "web_push"

func unmarshalData(raw []byte) map[string]any {
	if raw == nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	return m
}

type notificationRepo struct {
	db *db.DB
}

func NewNotificationRepo(database *db.DB) *notificationRepo {
	return &notificationRepo{db: database}
}

func (r *notificationRepo) Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	id := uuid.NewString()
	now := time.Now()

	if input.IdempotencyKey != "" {
		var existing domain.Notification
		var dataRaw []byte
		var idemKey sql.NullString
		err := r.db.QueryRowContext(ctx, `
			SELECT id, user_id, type, title, body, data, action_url, read_at, deleted_at, idempotency_key, created_at, updated_at
			FROM notifications WHERE user_id = ? AND idempotency_key = ?
			ORDER BY created_at DESC LIMIT 1
		`, input.RecipientID, input.IdempotencyKey).Scan(
			&existing.ID, &existing.UserID, &existing.Type, &existing.Title, &existing.Body, &dataRaw,
			&existing.ActionURL, &existing.ReadAt, &existing.DeletedAt, &idemKey, &existing.CreatedAt, &existing.UpdatedAt,
		)
		if err == nil {
			existing.Data = unmarshalData(dataRaw)
			existing.IdempotencyKey = idemKey.String
			return &existing, nil
		} else if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to check idempotency: %w", err)
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("notification: begin tx: %w", err)
	}
	defer tx.Rollback()

	dataJSON, _ := json.Marshal(input.Data)

	var idempotencyKey interface{} = input.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = nil
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (id, user_id, type, title, body, data, action_url, read_at, deleted_at, idempotency_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, input.RecipientID, input.Type, input.Title, input.Body, dataJSON, input.ActionURL, nil, nil, idempotencyKey, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	subs, err := NewPushSubscriptionRepo(r.db).ListByUser(ctx, input.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions for notification deliveries: %w", err)
	}

	for _, sub := range subs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO notification_deliveries (id, notification_id, subscription_id, channel, status, available_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'), datetime('now'))
		`, uuid.NewString(), id, sub.ID, DeliveryChannelWebPush, domain.DeliveryStatusPending)
		if err != nil {
			return nil, fmt.Errorf("failed to create delivery job for subscription: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit notification create: %w", err)
	}

	return r.Get(ctx, id, input.RecipientID)
}

func (r *notificationRepo) Get(ctx context.Context, id string, recipientID string) (*domain.Notification, error) {
	var n domain.Notification
	var dataRaw []byte
	var idemKey sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, type, title, body, data, action_url, read_at, deleted_at, idempotency_key, created_at, updated_at
		FROM notifications WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, id, recipientID).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &dataRaw, &n.ActionURL, &n.ReadAt, &n.DeletedAt, &idemKey, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	n.Data = unmarshalData(dataRaw)
	n.IdempotencyKey = idemKey.String
	return &n, nil
}

var _ domain.NotificationRepository = (*notificationRepo)(nil)

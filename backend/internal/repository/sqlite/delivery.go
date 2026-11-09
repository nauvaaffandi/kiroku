package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

type deliveryRepo struct{ db *db.DB }

func newDeliveryRepo(database *db.DB) *deliveryRepo { return &deliveryRepo{db: database} }

// NewDeliveryRepo returns a delivery repository implementing domain.DeliveryRepository.
func NewDeliveryRepo(database *db.DB) domain.DeliveryRepository {
	return &deliveryRepo{db: database}
}

func (r *deliveryRepo) Create(ctx context.Context, input domain.CreateDeliveryInput) (*domain.DeliveryJob, error) {
	id := uuid.NewString()
	avail := input.AvailableAt
	if avail.IsZero() {
		avail = time.Now()
	}
	availStr := avail.UTC().Format("2006-01-02 15:04:05")
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notification_deliveries (id, notification_id, subscription_id, channel, status, available_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, id, input.NotificationID, input.SubscriptionID, input.Channel, domain.DeliveryStatusPending, availStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create delivery job: %w", err)
	}
	return r.Get(ctx, id)
}

func scanDelivery(scanner interface{ Scan(dest ...any) error }) (*domain.DeliveryJob, error) {
	var j domain.DeliveryJob
	var lockToken, lastError sql.NullString
	if err := scanner.Scan(
		&j.ID, &j.NotificationID, &j.SubscriptionID, &j.Channel, &j.Status, &j.Attempts,
		&j.AvailableAt, &j.LockedAt, &lockToken, &j.SentAt, &lastError,
		&j.CreatedAt, &j.UpdatedAt,
	); err != nil {
		return nil, err
	}
	j.LockToken = lockToken.String
	j.LastError = lastError.String
	return &j, nil
}

func (r *deliveryRepo) Get(ctx context.Context, id string) (*domain.DeliveryJob, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, notification_id, subscription_id, channel, status, attempts, available_at, locked_at, lock_token, sent_at, last_error, created_at, updated_at
		FROM notification_deliveries WHERE id = ?
	`, id)
	j, err := scanDelivery(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	return j, nil
}

func (r *deliveryRepo) ClaimPending(ctx context.Context, limit int, lockTimeout time.Duration) ([]*domain.DeliveryJob, error) {
	token := uuid.NewString()
	minutes := int(lockTimeout.Minutes())
	if minutes <= 0 {
		minutes = 1
	}
	modifier := fmt.Sprintf("-%d minutes", minutes)

	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_deliveries
		SET status = ?, locked_at = datetime('now'), lock_token = ?, attempts = attempts + 1
		WHERE id IN (
			SELECT id FROM notification_deliveries
			WHERE (status = ? AND datetime(available_at) <= datetime('now'))
			   OR (status = ? AND datetime(locked_at) < datetime('now', ?))
			ORDER BY id
			LIMIT ?
		)
	`, domain.DeliveryStatusProcessing, token, domain.DeliveryStatusPending, domain.DeliveryStatusProcessing, modifier, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to claim deliveries: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, notification_id, subscription_id, channel, status, attempts, available_at, locked_at, lock_token, sent_at, last_error, created_at, updated_at
		FROM notification_deliveries WHERE lock_token = ?
	`, token)
	if err != nil {
		return nil, fmt.Errorf("failed to read claimed deliveries: %w", err)
	}
	defer rows.Close()

	var list []*domain.DeliveryJob
	for rows.Next() {
		j, err := scanDelivery(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan claimed delivery: %w", err)
		}
		list = append(list, j)
	}
	return list, rows.Err()
}

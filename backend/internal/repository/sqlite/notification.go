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

// Delivery channel constants reused by both repos.
const (
	DeliveryChannelWebPush = "web_push"
)

// unmarshalData converts JSON bytes from the data column into a map.
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

// --- NotificationRepo -------------------------------------------------------

type notificationRepo struct {
	db *db.DB
}

func NewNotificationRepo(database *db.DB) *notificationRepo {
	return &notificationRepo{db: database}
}

func (r *notificationRepo) Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	id := uuid.NewString()
	now := time.Now()

	// Idempotency: if key supplied and an identical logical event exists for this
	// recipient, return it without creating duplicates or delivery jobs.
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

	// Serialize input.Data as JSON for SQLite storage
	dataJSON, _ := json.Marshal(input.Data)

	var idempotencyKey interface{} = input.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = nil // store NULL so multiple no-key notifications don't collide
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

// GetByID returns a notification by ID without an owner scope (worker-facing).
func (r *notificationRepo) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	var n domain.Notification
	var dataRaw []byte
	var idemKey sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, type, title, body, data, action_url, read_at, deleted_at, idempotency_key, created_at, updated_at
		FROM notifications WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &dataRaw, &n.ActionURL, &n.ReadAt, &n.DeletedAt, &idemKey, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification by id: %w", err)
	}
	n.Data = unmarshalData(dataRaw)
	n.IdempotencyKey = idemKey.String
	return &n, nil
}

func (r *notificationRepo) List(ctx context.Context, recipientID string, filter domain.ListFilter) ([]*domain.Notification, error) {
	q := `
		SELECT id, user_id, type, title, body, data, action_url, read_at, deleted_at, idempotency_key, created_at, updated_at
		FROM notifications WHERE user_id = ? AND deleted_at IS NULL
	`
	args := []interface{}{recipientID}

	if filter.UnreadOnly {
		q += ` AND read_at IS NULL`
	}

	if filter.Limit > 0 {
		q += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
		args = append(args, filter.Limit, filter.Offset)
	} else {
		q += ` ORDER BY created_at DESC`
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var list []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		var dataRaw []byte
		var idemKey sql.NullString
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &dataRaw, &n.ActionURL, &n.ReadAt, &n.DeletedAt, &idemKey, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		n.Data = unmarshalData(dataRaw)
		n.IdempotencyKey = idemKey.String
		list = append(list, &n)
	}
	if list == nil {
		list = []*domain.Notification{}
	}
	return list, rows.Err()
}

func (r *notificationRepo) Count(ctx context.Context, recipientID string, unreadOnly bool) (int, error) {
	q := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND deleted_at IS NULL`
	args := []interface{}{recipientID}
	if unreadOnly {
		q += ` AND read_at IS NULL`
	}

	var n int
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}
	return n, nil
}

func (r *notificationRepo) MarkRead(ctx context.Context, id string, recipientID string) (*domain.Notification, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?
	`, id, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification read: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return nil, sql.ErrNoRows
	}
	return r.Get(ctx, id, recipientID)
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, recipientID string) (int, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE user_id = ? AND read_at IS NULL AND deleted_at IS NULL
	`, recipientID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark all notifications read: %w", err)
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func (r *notificationRepo) Delete(ctx context.Context, id string, recipientID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("notification delete: begin tx: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `UPDATE notifications SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ? AND deleted_at IS NULL`, id, recipientID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM notification_deliveries WHERE notification_id = ? AND status IN (?, ?)`, id, domain.DeliveryStatusPending, domain.DeliveryStatusProcessing); err != nil {
		return fmt.Errorf("failed to cancel deliveries: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit notification delete: %w", err)
	}
	return nil
}

func (r *notificationRepo) DeleteRead(ctx context.Context, recipientID string) (int, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM notifications WHERE user_id = ? AND read_at IS NOT NULL AND deleted_at IS NULL
	`, recipientID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete read notifications: %w", err)
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

var _ domain.NotificationRepository = (*notificationRepo)(nil)

// --- DeliveryRepo -----------------------------------------------------------

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

// ClaimPending atomically claims up to limit due jobs that are pending or
// whose processing lock has expired (crash recovery). Returns claimed jobs.
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

func (r *deliveryRepo) MarkSent(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_deliveries SET status = ?, sent_at = datetime('now'), last_error = '', locked_at = NULL, lock_token = NULL WHERE id = ?
	`, domain.DeliveryStatusSent, id)
	if err != nil {
		return fmt.Errorf("failed to mark delivery sent: %w", err)
	}
	return nil
}

func (r *deliveryRepo) MarkFailed(ctx context.Context, id, lastError string) error {
	return r.updateStatus(ctx, id, domain.DeliveryStatusFailed, lastError)
}

func (r *deliveryRepo) Reschedule(ctx context.Context, id string, availableAt time.Time, lastError string) error {
	availStr := availableAt.UTC().Format("2006-01-02 15:04:05")
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_deliveries SET status = ?, available_at = ?, locked_at = NULL, lock_token = NULL, last_error = ? WHERE id = ?
	`, domain.DeliveryStatusPending, availStr, lastError, id)
	return err
}

func (r *deliveryRepo) CancelByNotification(ctx context.Context, notificationID string) (int, error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM notification_deliveries WHERE notification_id = ? AND status IN (?, ?)
	`, notificationID, domain.DeliveryStatusPending, domain.DeliveryStatusProcessing)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel deliveries: %w", err)
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func (r *deliveryRepo) updateStatus(ctx context.Context, id, status, lastError string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE notification_deliveries SET status = ?, last_error = ?, locked_at = NULL, lock_token = NULL WHERE id = ?
	`, status, lastError, id)
	if err != nil {
		return fmt.Errorf("failed to update delivery status: %w", err)
	}
	return nil
}

var _ domain.DeliveryRepository = (*deliveryRepo)(nil)

// scanNotificationData converts raw JSON bytes to map[string]any.
func scanData(b []byte, target *map[string]any) error {
	if b == nil {
		return nil
	}
	return json.Unmarshal(b, target)
}

// var _ strings for compile check removed

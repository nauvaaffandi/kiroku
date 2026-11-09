package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"kiroku/backend/internal/domain"
)

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

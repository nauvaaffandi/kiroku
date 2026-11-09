package sqlite

import (
	"context"

	"encoding/json"
	"fmt"
	"time"

	"kiroku/backend/internal/domain"
)

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

func scanData(b []byte, target *map[string]any) error {
	if b == nil {
		return nil
	}
	return json.Unmarshal(b, target)
}

var _ domain.DeliveryRepository = (*deliveryRepo)(nil)

package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"kiroku/backend/internal/domain"
)

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

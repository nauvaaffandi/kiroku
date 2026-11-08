package domain

import (
	"context"
	"time"
)

// Notification is an application record in the user's notification inbox.
// It is domain-agnostic: type and data are extensible by any business module.
type Notification struct {
	ID             string         `json:"id"`
	UserID         string         `json:"user_id"`
	Type           string         `json:"type"`
	Title          string         `json:"title"`
	Body           string         `json:"body"`
	Data           map[string]any `json:"data,omitempty"`
	ActionURL      string         `json:"action_url,omitempty"`
	ReadAt         *time.Time     `json:"read_at,omitempty"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// CreateNotificationInput holds everything needed to create a notification.
type CreateNotificationInput struct {
	RecipientID    string
	Type           string
	Title          string
	Body           string
	Data           map[string]any
	ActionURL      string
	IdempotencyKey string
}

// ListFilter controls how notifications are returned by List.
type ListFilter struct {
	UnreadOnly bool
	Limit      int
	Offset     int
}

// NotificationRepository persists notification records and manages their state.
// Implementations must not depend on business-specific modules.
type NotificationRepository interface {
	// Create stores a notification and returns it. The supplied IdempotencyKey,
	// when non-empty, prevents duplicate records for the same logical event.
	Create(ctx context.Context, input CreateNotificationInput) (*Notification, error)
// Get returns a single notification owned by recipientID, or nil if not found.
Get(ctx context.Context, id, recipientID string) (*Notification, error)
// GetByID returns a single notification by ID regardless of owner (used by worker).
GetByID(ctx context.Context, id string) (*Notification, error)
	// List returns notifications for recipientID, newest first.
	List(ctx context.Context, recipientID string, filter ListFilter) ([]*Notification, error)
	// Count returns the number of notifications for recipientID.
	Count(ctx context.Context, recipientID string, unreadOnly bool) (int, error)
	// MarkRead marks a single notification read and returns it.
	MarkRead(ctx context.Context, id, recipientID string) (*Notification, error)
	// MarkAllRead marks every unread notification read, returns how many changed.
	MarkAllRead(ctx context.Context, recipientID string) (int, error)
	// Delete soft-deletes a notification owned by recipientID.
	Delete(ctx context.Context, id, recipientID string) error
	// DeleteRead removes all read notifications owned by recipientID.
	DeleteRead(ctx context.Context, recipientID string) (int, error)
}

// DeliveryJob is a persistent push-delivery attempt. One notification can have
// many delivery jobs — one per subscribed browser/device. Jobs outlive delivery:
// they record success and failure history.
type DeliveryJob struct {
	ID             string     `json:"id"`
	NotificationID string     `json:"notification_id"`
	SubscriptionID string     `json:"subscription_id"`
	Channel        string     `json:"channel"`
	Status         string     `json:"status"`
	Attempts       int        `json:"attempts"`
	AvailableAt    time.Time  `json:"available_at"`
	LockedAt       *time.Time `json:"locked_at,omitempty"`
	LockToken      string     `json:"-"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	LastError      string     `json:"last_error,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

const (
	DeliveryChannelWebPush = "web_push"

	DeliveryStatusPending    = "pending"
	DeliveryStatusProcessing = "processing"
	DeliveryStatusSent       = "sent"
	DeliveryStatusFailed     = "failed"
	DeliveryStatusCancelled  = "cancelled"
)

type CreateDeliveryInput struct {
	NotificationID string
	SubscriptionID string
	Channel        string
	AvailableAt    time.Time
}

// DeliveryRepository persists delivery jobs and the claim/lock lifecycle.
// Deliveries are the source of truth for pending push attempts.
type DeliveryRepository interface {
	// Create inserts a delivery job and returns it.
	Create(ctx context.Context, input CreateDeliveryInput) (*DeliveryJob, error)
	// Get returns a delivery job by ID.
	Get(ctx context.Context, id string) (*DeliveryJob, error)
	// ClaimPending atomically marks up to limit due jobs as processing and returns them.
	ClaimPending(ctx context.Context, limit int, lockTimeout time.Duration) ([]*DeliveryJob, error)
	// MarkSent records a successful delivery.
	MarkSent(ctx context.Context, id string) error
	// MarkFailed records a permanent failure.
	MarkFailed(ctx context.Context, id, lastError string) error
	// Reschedule moves a job back to pending at availableAt with a recorded error.
	Reschedule(ctx context.Context, id string, availableAt time.Time, lastError string) error
	// CancelByNotification cancels all non-terminal jobs for a deleted notification.
	CancelByNotification(ctx context.Context, notificationID string) (int, error)
}

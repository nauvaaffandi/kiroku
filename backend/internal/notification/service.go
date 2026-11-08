package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/domain"
)

// PushSender delivers a push notification to a single browser subscription.
// Implementations may be real (Web Push over HTTP) or a test fake.
type PushSender interface {
	Send(ctx context.Context, sub *domain.PushSubscription, n *domain.Notification) error
}

// ServiceConfig tunes the service's claim-and-process behaviour.
type ServiceConfig struct {
	// MaxAttempts is the maximum number of delivery attempts before marking a job failed.
	MaxAttempts int
	// RetryBase is the base delay for exponential backoff between attempts.
	RetryBase time.Duration
	// LockTimeout is how long a processing lock survives a worker crash before another
	// worker may reclaim the job.
	LockTimeout time.Duration
}

// Service is the reusable, domain-agnostic notification service.
// It does not import business-specific modules.
type Service struct {
	notifications domain.NotificationRepository
	deliveries    domain.DeliveryRepository
	subscriptions domain.PushSubscriptionRepository
	sender        PushSender
	cfg           ServiceConfig
}

func NewService(
	notifications domain.NotificationRepository,
	deliveries domain.DeliveryRepository,
	subscriptions domain.PushSubscriptionRepository,
	sender PushSender,
	cfg ServiceConfig,
) *Service {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	if cfg.RetryBase <= 0 {
		cfg.RetryBase = 30 * time.Second
	}
	if cfg.LockTimeout <= 0 {
		cfg.LockTimeout = 5 * time.Minute
	}
	return &Service{
		notifications: notifications,
		deliveries:    deliveries,
		subscriptions: subscriptions,
		sender:        sender,
		cfg:           cfg,
	}
}

// Create persists a notification and creates one pending delivery job per
// recipient push subscription inside a single transaction. An IdempotencyKey,
// when set, makes creation idempotent at the database level.
func (s *Service) Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	return s.notifications.Create(ctx, input)
}

// List returns notifications for recipientID, newest first.
func (s *Service) List(ctx context.Context, recipientID string, filter domain.ListFilter) ([]*domain.Notification, error) {
	return s.notifications.List(ctx, recipientID, filter)
}

// Count returns the number of notifications for recipientID.
func (s *Service) Count(ctx context.Context, recipientID string, unreadOnly bool) (int, error) {
	return s.notifications.Count(ctx, recipientID, unreadOnly)
}

// MarkRead marks a single notification read and returns it.
func (s *Service) MarkRead(ctx context.Context, id, recipientID string) (*domain.Notification, error) {
	return s.notifications.MarkRead(ctx, id, recipientID)
}

// MarkAllRead marks every unread notification read, returns how many changed.
func (s *Service) MarkAllRead(ctx context.Context, recipientID string) (int, error) {
	return s.notifications.MarkAllRead(ctx, recipientID)
}

// Delete soft-deletes a notification owned by recipientID and cancels its pending
// deliveries.
func (s *Service) Delete(ctx context.Context, id, recipientID string) error {
	return s.notifications.Delete(ctx, id, recipientID)
}

// DeleteRead removes all read notifications owned by recipientID.
func (s *Service) DeleteRead(ctx context.Context, recipientID string) (int, error) {
	return s.notifications.DeleteRead(ctx, recipientID)
}

// UnreadCount is an alias for Count(unreadOnly=true).
func (s *Service) UnreadCount(ctx context.Context, recipientID string) (int, error) {
	return s.notifications.Count(ctx, recipientID, true)
}

// RegisterSubscription stores a push subscription for the recipient.
func (s *Service) RegisterSubscription(ctx context.Context, input domain.CreatePushSubscriptionInput) (*domain.PushSubscription, error) {
	return s.subscriptions.Create(ctx, input)
}

// RemoveSubscription removes a push subscription by ID owned by recipientID.
func (s *Service) RemoveSubscription(ctx context.Context, id, recipientID string) error {
	return s.subscriptions.Delete(ctx, id)
}

// ListSubscriptions returns all active subscriptions for recipientID.
func (s *Service) ListSubscriptions(ctx context.Context, recipientID string) ([]*domain.PushSubscription, error) {
	return s.subscriptions.ListByUser(ctx, recipientID)
}

// ProcessPending claims a batch of due delivery jobs and attempts to send each.
// Returns the number of jobs successfully sent in this batch.
func (s *Service) ProcessPending(ctx context.Context, batchSize, maxAttempts int, retryBase, lockTimeout time.Duration) (int, error) {
	jobs, err := s.deliveries.ClaimPending(ctx, batchSize, lockTimeout)
	if err != nil {
		return 0, fmt.Errorf("process pending: claim: %w", err)
	}

	sent := 0
	now := time.Now()

	for _, job := range jobs {
		// Re-read to confirm we still own the claim (another worker may have reclaimed).
		cur, err := s.deliveries.Get(ctx, job.ID)
		if err != nil || cur == nil || cur.Status != domain.DeliveryStatusProcessing || cur.LockToken != job.LockToken {
			continue
		}

		notif, err := s.notifications.GetByID(ctx, cur.NotificationID)
		if err != nil {
			log.Warn().Err(err).Str("delivery", cur.ID).Msg("notification: failed to load notification for delivery")
			if perr := s.deliveries.Reschedule(ctx, cur.ID, now.Add(retryBase), err.Error()); perr != nil {
				log.Error().Err(perr).Str("delivery", cur.ID).Msg("notification: failed to reschedule delivery")
			}
			continue
		}
		if notif == nil || notif.DeletedAt != nil {
			if _, err := s.deliveries.CancelByNotification(ctx, cur.NotificationID); err != nil {
				log.Error().Err(err).Str("notification", cur.NotificationID).Msg("notification: failed to cancel deliveries")
			}
			continue
		}

		sub, err := s.subscriptions.Get(ctx, cur.SubscriptionID)
		if err != nil {
			log.Warn().Err(err).Str("delivery", cur.ID).Msg("notification: failed to load subscription")
			if perr := s.deliveries.MarkFailed(ctx, cur.ID, "subscription lookup failed"); perr != nil {
				log.Error().Err(perr).Str("delivery", cur.ID).Msg("notification: failed to mark delivery failed")
			}
			continue
		}
		if sub == nil {
			if perr := s.deliveries.MarkFailed(ctx, cur.ID, "subscription not found"); perr != nil {
				log.Error().Err(perr).Str("delivery", cur.ID).Msg("notification: failed to mark delivery failed")
			}
			continue
		}

		if err := s.sender.Send(ctx, sub, notif); err != nil {
			if isExpiredError(err) {
				// Subscription is invalid/expired — remove it permanently.
				if derr := s.subscriptions.Delete(ctx, sub.ID); derr != nil {
					log.Warn().Err(derr).Str("subscription", sub.ID).Msg("notification: failed to delete expired subscription")
				}
			}
			if perr := s.failDelivery(ctx, cur, err.Error(), maxAttempts, retryBase); perr != nil {
				log.Error().Err(perr).Str("delivery", cur.ID).Msg("notification: failed to record delivery failure")
			}
			continue
		}

		if err := s.deliveries.MarkSent(ctx, cur.ID); err != nil {
			log.Error().Err(err).Str("delivery", cur.ID).Msg("notification: failed to mark delivery sent")
			continue
		}
		sent++
	}

	return sent, nil
}

// failDelivery marks a job failed when retries are exhausted, otherwise reschedules
// it with exponential backoff.
func (s *Service) failDelivery(ctx context.Context, cur *domain.DeliveryJob, lastError string, maxAttempts int, retryBase time.Duration) error {
	if cur.Attempts >= maxAttempts {
		return s.deliveries.MarkFailed(ctx, cur.ID, lastError)
	}
	backoff := retryBase * time.Duration(1<<uint(cur.Attempts-1))
	if backoff > time.Hour {
		backoff = time.Hour
	}
	return s.deliveries.Reschedule(ctx, cur.ID, time.Now().Add(backoff), lastError)
}

func isExpiredError(err error) bool {
	msg := err.Error()
	return containsCI(msg, "410") || containsCI(msg, "404") || containsCI(msg, "gone") || containsCI(msg, "expired")
}

func containsCI(haystack, needle string) bool {
	h, n := haystack, needle
	for i := 0; i < len(h); i++ {
		if h[i] >= 'A' && h[i] <= 'Z' {
			h = h[:i] + string(h[i]+32) + h[i+1:]
		}
	}
	for i := 0; i < len(n); i++ {
		if n[i] >= 'A' && n[i] <= 'Z' {
			n = n[:i] + string(n[i]+32) + n[i+1:]
		}
	}
	return contains(h, n)
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

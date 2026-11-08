package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

func newNotificationTestEnv(t *testing.T) (*db.DB, *notificationRepo, domain.DeliveryRepository, *PushSubscriptionRepo) {
	t.Helper()
	dataDir := t.TempDir() + "/data"
	d, err := db.Open(dataDir)
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO users (id, email, pin_hash) VALUES (?, ?, ?)`, "user_1", "user_1@example.com", "hash")
	require.NoError(t, err)

	return d, NewNotificationRepo(d), NewDeliveryRepo(d), NewPushSubscriptionRepo(d)
}

func TestNotificationCreateListAndUnreadCount(t *testing.T) {
	d, notifRepo, _, _ := newNotificationTestEnv(t)
	defer d.Close()

	n, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{
		RecipientID: "user_1",
		Type:        "info",
		Title:       "Welcome",
		Body:        "Hello world",
		Data:        map[string]any{"app": "kiroku"},
	})
	require.NoError(t, err)
	require.NotNil(t, n)
	require.NotEmpty(t, n.ID)
	assert.Equal(t, "user_1", n.UserID)
	assert.Equal(t, "info", n.Type)
	assert.Equal(t, "Welcome", n.Title)
	assert.Equal(t, "Hello world", n.Body)
	assert.Equal(t, "kiroku", n.Data["app"])

	list, err := notifRepo.List(context.Background(), "user_1", domain.ListFilter{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.False(t, list[0].ReadAt != nil)

	count, err := notifRepo.Count(context.Background(), "user_1", true)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	countRead, err := notifRepo.Count(context.Background(), "user_1", false)
	require.NoError(t, err)
	assert.Equal(t, 1, countRead)
}

func TestNotificationMarkReadAndMarkAllRead(t *testing.T) {
	d, notifRepo, _, _ := newNotificationTestEnv(t)
	defer d.Close()

	n1, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "One", Body: "one body"})
	require.NoError(t, err)
	_, err = notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "Two", Body: "two body"})
	require.NoError(t, err)

	// Mark one
	updated, err := notifRepo.MarkRead(context.Background(), n1.ID, "user_1")
	require.NoError(t, err)
	assert.NotNil(t, updated.ReadAt)

	count, err := notifRepo.Count(context.Background(), "user_1", true)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Mark all
	n, err := notifRepo.MarkAllRead(context.Background(), "user_1")
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	count, err = notifRepo.Count(context.Background(), "user_1", true)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Cross-user IDOR: user_2 cannot mark user_1's notification
	_, err = notifRepo.MarkRead(context.Background(), n1.ID, "user_2")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestNotificationIdempotentCreation(t *testing.T) {
	d, notifRepo, _, _ := newNotificationTestEnv(t)
	defer d.Close()

	key := uuid.NewString()
	n1, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{
		RecipientID:   "user_1",
		Type:          "info",
		Title:         "Same",
		Body:          "body",
		IdempotencyKey: key,
	})
	require.NoError(t, err)

	n2, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{
		RecipientID:   "user_1",
		Type:          "info",
		Title:         "Same retry",
		Body:          "different body",
		IdempotencyKey: key,
	})
	require.NoError(t, err)
	assert.Equal(t, n1.ID, n2.ID)
	assert.Equal(t, "body", n2.Body)

	// Count should be 1 (no duplicate)
	count, err := notifRepo.Count(context.Background(), "user_1", false)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestNotificationDeleteCancelsPendingDeliveries(t *testing.T) {
	d, notifRepo, deliveryRepo, pushRepo := newNotificationTestEnv(t)
	defer d.Close()

	// Create a subscription for user_1
	sub, err := pushRepo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID: "user_1",
		Endpoint: "https://example.com/push/sub1",
		Keys: domain.Keys{P256DH: "p256dh", Auth: "auth"},
		UserAgent: "test",
	})
	require.NoError(t, err)

	n, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{
		RecipientID: "user_1",
		Type:          "info",
		Title:         "Deliverable",
		Body:          "body",
	})
	require.NoError(t, err)

	// Delivery job should exist and become processing upon claim
	deliveries, err := deliveryRepo.ClaimPending(context.Background(), 10, time.Minute)
	require.NoError(t, err)
	require.Len(t, deliveries, 1)
	assert.Equal(t, sub.ID, deliveries[0].SubscriptionID)
	assert.Equal(t, n.ID, deliveries[0].NotificationID)
	assert.Equal(t, domain.DeliveryStatusProcessing, deliveries[0].Status)

	// Delete notification — soft delete + cancel pending deliveries
	err = notifRepo.Delete(context.Background(), n.ID, "user_1")
	require.NoError(t, err)

	// Notification gone
	gone, err := notifRepo.Get(context.Background(), n.ID, "user_1")
	require.NoError(t, err)
	assert.Nil(t, gone)

	// Pending delivery cancelled
	remaining, err := deliveryRepo.ClaimPending(context.Background(), 10, time.Minute)
	require.NoError(t, err)
	assert.Len(t, remaining, 0)
}

func TestNotificationDeleteRead(t *testing.T) {
	d, notifRepo, _, _ := newNotificationTestEnv(t)
	defer d.Close()

	n1, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "One", Body: "one"})
	require.NoError(t, err)
	_, err = notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "Two", Body: "two"})
	require.NoError(t, err)

	_, err = notifRepo.MarkRead(context.Background(), n1.ID, "user_1")
	require.NoError(t, err)

	deleted, err := notifRepo.DeleteRead(context.Background(), "user_1")
	require.NoError(t, err)
	assert.Equal(t, 1, deleted)

	count, err := notifRepo.Count(context.Background(), "user_1", false)
	require.NoError(t, err)
	assert.Equal(t, 1, count) // unread one remains
}

func TestNotificationListUnreadOnly(t *testing.T) {
	d, notifRepo, _, _ := newNotificationTestEnv(t)
	defer d.Close()

	_, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "One", Body: "one"})
	require.NoError(t, err)
	_, err = notifRepo.Create(context.Background(), domain.CreateNotificationInput{RecipientID: "user_1", Type: "info", Title: "Two", Body: "two"})
	require.NoError(t, err)

	all, err := notifRepo.List(context.Background(), "user_1", domain.ListFilter{Limit: 10})
	require.NoError(t, err)
	assert.Len(t, all, 2)

	unread, err := notifRepo.List(context.Background(), "user_1", domain.ListFilter{UnreadOnly: true, Limit: 10})
	require.NoError(t, err)
	assert.Len(t, unread, 2)
}

func TestNotificationDeliveryJobRetryAndSent(t *testing.T) {
	d, notifRepo, deliveryRepo, pushRepo := newNotificationTestEnv(t)
	defer d.Close()

	sub, err := pushRepo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID:    "user_1",
		Endpoint:  "https://example.com/push/job_test",
		Keys:      domain.Keys{P256DH: "k", Auth: "a"},
		UserAgent: "agent",
	})
	require.NoError(t, err)

	notif, err := notifRepo.Create(context.Background(), domain.CreateNotificationInput{
		RecipientID: "user_1",
		Type:        "info",
		Title:       "Test",
		Body:        "Body",
	})
	require.NoError(t, err)

	// A job was created automatically by notifRepo.Create. Clear it first to test manual Create:
	_, err = d.ExecContext(context.Background(), "DELETE FROM notification_deliveries")
	require.NoError(t, err)

	now := time.Now()
	job, err := deliveryRepo.Create(context.Background(), domain.CreateDeliveryInput{
		NotificationID: notif.ID,
		SubscriptionID: sub.ID,
		Channel:        domain.DeliveryChannelWebPush,
		AvailableAt:    now,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.DeliveryStatusPending, job.Status)
	assert.Equal(t, 0, job.Attempts)

	// Claim
	claimed, err := deliveryRepo.ClaimPending(context.Background(), 10, time.Minute)
	require.NoError(t, err)
	assert.Len(t, claimed, 1)
	assert.Equal(t, domain.DeliveryStatusProcessing, claimed[0].Status)
	assert.Equal(t, 1, claimed[0].Attempts)
	require.NotEmpty(t, claimed[0].LockToken)

	// Mark failed -> rescheduled with backoff
	require.NoError(t, deliveryRepo.Reschedule(context.Background(), claimed[0].ID, now.Add(time.Minute), "network blip"))
	got, err := deliveryRepo.Get(context.Background(), claimed[0].ID)
	require.NoError(t, err)
	assert.Equal(t, domain.DeliveryStatusPending, got.Status)
	assert.True(t, got.AvailableAt.After(now))

	// Mark sent
	require.NoError(t, deliveryRepo.MarkSent(context.Background(), claimed[0].ID))
	got, err = deliveryRepo.Get(context.Background(), claimed[0].ID)
	require.NoError(t, err)
	assert.Equal(t, domain.DeliveryStatusSent, got.Status)
	assert.NotNil(t, got.SentAt)
}

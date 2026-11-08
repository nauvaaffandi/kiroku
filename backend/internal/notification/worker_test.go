package notification_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/notification"
	"kiroku/backend/internal/repository/sqlite"
)

type noopSender struct{}

func (s *noopSender) Send(ctx context.Context, sub *domain.PushSubscription, n *domain.Notification) error {
	return nil
}

func setupTestService(t *testing.T) (*notification.Service, *db.DB) {
	t.Helper()
	dataDir := t.TempDir() + "/data"
	d, err := db.Open(dataDir)
	require.NoError(t, err)

	_, err = d.ExecContext(context.Background(), `INSERT INTO users (id, email, pin_hash) VALUES ('u1', 'u1@test.com', 'h')`)
	require.NoError(t, err)

	notifRepo := sqlite.NewNotificationRepo(d)
	deliveryRepo := sqlite.NewDeliveryRepo(d)
	pushRepo := sqlite.NewPushSubscriptionRepo(d)

	svc := notification.NewService(notifRepo, deliveryRepo, pushRepo, &noopSender{}, notification.ServiceConfig{
		MaxAttempts: 3,
		RetryBase:   10 * time.Millisecond,
		LockTimeout: time.Minute,
	})
	return svc, d
}

func TestWorker_DisabledStopNoPanic(t *testing.T) {
	svc, d := setupTestService(t)
	defer d.Close()

	w := notification.NewWorker(svc, notification.Config{
		Enabled: false,
	}, context.Background())

	// Stop must not panic when worker is disabled
	assert.NotPanics(t, func() {
		w.Stop()
	})
}

func TestWorker_StartAndStopCleanly(t *testing.T) {
	svc, d := setupTestService(t)
	defer d.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := notification.NewWorker(svc, notification.Config{
		Interval:    10 * time.Millisecond,
		BatchSize:   10,
		MaxAttempts: 3,
		RetryBase:   10 * time.Millisecond,
		LockTimeout: time.Minute,
		Enabled:     true,
	}, ctx)

	time.Sleep(30 * time.Millisecond)

	assert.NotPanics(t, func() {
		w.Stop()
	})
}

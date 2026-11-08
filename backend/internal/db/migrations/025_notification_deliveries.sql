-- Migration 025: Notification delivery jobs (persistent push delivery queue)
CREATE TABLE IF NOT EXISTS notification_deliveries (
    id TEXT PRIMARY KEY,
    notification_id TEXT NOT NULL,
    subscription_id TEXT NOT NULL,
    channel TEXT NOT NULL DEFAULT 'web_push',
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    lock_token TEXT,
    available_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    locked_at DATETIME,
    sent_at DATETIME,
    last_error TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES push_subscriptions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_deliveries_pending ON notification_deliveries(status, available_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_deliveries_notification ON notification_deliveries(notification_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_subscription ON notification_deliveries(subscription_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_lock_token ON notification_deliveries(lock_token);

CREATE TRIGGER IF NOT EXISTS notification_deliveries_updated_at
AFTER UPDATE ON notification_deliveries
BEGIN
    UPDATE notification_deliveries SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

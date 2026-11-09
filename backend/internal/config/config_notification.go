package config

import (
	"os"
	"strconv"
	"time"
)

// NotificationWorkerInterval returns the notification delivery worker interval.
// Notification delivery worker interval (default 30s).
func NotificationWorkerInterval() time.Duration {
	if v := os.Getenv("NOTIFICATION_WORKER_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 30 * time.Second
}

// NotificationWorkerBatchSize returns the number of delivery jobs processed per pass.
// Default 10.
func NotificationWorkerBatchSize() int {
	if v := os.Getenv("NOTIFICATION_WORKER_BATCH_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 10
}

// NotificationWorkerMaxAttempts returns the maximum delivery attempts per job before
// marking it as failed. Default 5.
func NotificationWorkerMaxAttempts() int {
	if v := os.Getenv("NOTIFICATION_MAX_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 5
}

// NotificationWorkerRetryBase returns the base delay for exponential backoff between
// delivery attempts. Default 30s.
func NotificationWorkerRetryBase() time.Duration {
	if v := os.Getenv("NOTIFICATION_RETRY_BASE_DELAY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 30 * time.Second
}

// NotificationWorkerLockTimeout returns how long a processing lock survives a worker
// crash before another worker may reclaim the job. Default 5m.
func NotificationWorkerLockTimeout() time.Duration {
	if v := os.Getenv("NOTIFICATION_LOCK_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 5 * time.Minute
}

// NotificationWorkerEnabled returns whether the notification delivery worker runs.
// Default true.
func NotificationWorkerEnabled() bool {
	return os.Getenv("NOTIFICATION_WORKER_ENABLED") != "false"
}

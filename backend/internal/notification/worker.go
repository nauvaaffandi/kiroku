// package notification provides a background worker that periodically processes pending
// push notification delivery jobs.
package notification

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Worker periodically runs to process due delivery jobs. It respects context cancellation
// for clean shutdown.
type Worker struct {
	service     *Service
	interval    time.Duration
	stop        chan struct{}
	stopped     chan struct{}
	mu          sync.Mutex
	running     bool
	cancel      context.CancelFunc
}

// Config describes settings for the notification delivery worker.
type Config struct {
	Interval     time.Duration
	BatchSize    int
	MaxAttempts  int
	RetryBase    time.Duration
	LockTimeout  time.Duration
	Enabled      bool
}

// NewWorker builds a notification delivery worker. The given context is used for
// cancellation; the worker listens to its Done() channel.
func NewWorker(service *Service, cfg Config, ctx context.Context) *Worker {
	w := &Worker{
		service:   service,
		interval:  cfg.Interval,
		stop:      make(chan struct{}),
		stopped:   make(chan struct{}),
		mu:        sync.Mutex{},
	}

	if !cfg.Enabled {
		// Disabled workers never start.
		close(w.stopped)
		return w
	}

	workerCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	go w.loop(workerCtx, cfg)
	return w
}

// Stop gracefully shuts down the worker; it blocks until the worker has stopped.
func (w *Worker) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
	}
	w.mu.Unlock()
	<-w.stopped
}

// loop runs the worker's main loop.
func (w *Worker) loop(ctx context.Context, cfg Config) {
	defer func() { close(w.stopped) }()

	w.mu.Lock()
	w.running = true
	w.mu.Unlock()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	// Initial immediate run so notifications are delivered promptly after creation.
	if _, err := w.service.ProcessPending(ctx, cfg.BatchSize, cfg.MaxAttempts, cfg.RetryBase, cfg.LockTimeout); err != nil {
		log.Error().Err(err).Msg("notification: failed to run initial delivery pass")
	}

	for {
		select {
		case <-ctx.Done():
			w.mu.Lock()
			w.running = false
			w.mu.Unlock()
			log.Debug().Msg("notification: worker exiting")
			return
		case <-ticker.C:
			if _, err := w.service.ProcessPending(ctx, cfg.BatchSize, cfg.MaxAttempts, cfg.RetryBase, cfg.LockTimeout); err != nil {
				log.Warn().Err(err).Msg("notification: delivery pass failed")
			}
		}
	}
}

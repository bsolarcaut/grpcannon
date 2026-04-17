// Package ticker provides a periodic tick utility that fires at a fixed
// interval and can be stopped cleanly via context cancellation.
package ticker

import (
	"context"
	"time"
)

// Ticker fires a callback at a fixed interval until the context is done
// or Stop is called.
type Ticker struct {
	interval time.Duration
	callback func()
	stop     chan struct{}
}

// New creates a Ticker that calls fn every interval.
// If interval is <= 0 it defaults to 1 second.
func New(interval time.Duration, fn func()) *Ticker {
	if interval <= 0 {
		interval = time.Second
	}
	return &Ticker{
		interval: interval,
		callback: fn,
		stop:     make(chan struct{}),
	}
}

// Start begins ticking in a new goroutine. It returns immediately.
// The ticker stops when ctx is cancelled or Stop is called.
func (t *Ticker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.callback()
			case <-ctx.Done():
				return
			case <-t.stop:
				return
			}
		}
	}()
}

// Stop halts the ticker. Safe to call multiple times.
func (t *Ticker) Stop() {
	select {
	case <-t.stop:
		// already stopped
	default:
		close(t.stop)
	}
}

// Interval returns the configured tick interval.
func (t *Ticker) Interval() time.Duration {
	return t.interval
}

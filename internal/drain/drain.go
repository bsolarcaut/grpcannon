// Package drain provides graceful shutdown helpers that wait for
// in-flight requests to complete before the process exits.
package drain

import (
	"context"
	"sync"
	"time"
)

// Drainer tracks in-flight work and blocks until all work is done or the
// context deadline is exceeded.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
	timeout time.Duration
}

// New returns a Drainer with the given drain timeout.
// If timeout is zero, DefaultTimeout is used.
const DefaultTimeout = 5 * time.Second

func New(timeout time.Duration) *Drainer {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Drainer{timeout: timeout}
}

// Acquire marks the start of a unit of work.
// Returns false if the Drainer has already been closed.
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks the completion of a unit of work.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Drain closes the Drainer to new acquisitions and waits for all in-flight
// work to finish, or until ctx is cancelled.
func (d *Drainer) Drain(ctx context.Context) error {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

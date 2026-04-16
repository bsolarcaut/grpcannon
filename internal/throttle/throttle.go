// Package throttle provides a concurrency limiter using a semaphore pattern.
package throttle

import (
	"context"
	"errors"
)

// ErrThrottled is returned when the throttle is full and the context is done.
var ErrThrottled = errors.New("throttle: context cancelled while waiting")

// Throttle limits the number of concurrent operations.
type Throttle struct {
	sem chan struct{}
}

// New creates a Throttle allowing up to n concurrent operations.
// If n <= 0, it defaults to 1.
func New(n int) *Throttle {
	if n <= 0 {
		n = 1
	}
	return &Throttle{sem: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is done.
func (t *Throttle) Acquire(ctx context.Context) error {
	select {
	case t.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrThrottled
	}
}

// Release frees a previously acquired slot.
func (t *Throttle) Release() {
	select {
	case <-t.sem:
	default:
	}
}

// Cap returns the maximum concurrency allowed.
func (t *Throttle) Cap() int {
	return cap(t.sem)
}

// Len returns the number of currently active slots.
func (t *Throttle) Len() int {
	return len(t.sem)
}

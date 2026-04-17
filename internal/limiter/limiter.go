// Package limiter provides a token-bucket style concurrency limiter
// that caps the number of in-flight gRPC calls at any given time.
package limiter

import (
	"context"
	"sync/atomic"
)

// Limiter gates concurrent execution up to a fixed capacity.
type Limiter struct {
	ch      chan struct{}
	acquired atomic.Int64
}

// New creates a Limiter with the given capacity.
// If capacity <= 0 it defaults to 1.
func New(capacity int) *Limiter {
	if capacity <= 0 {
		capacity = 1
	}
	return &Limiter{ch: make(chan struct{}, capacity)}
}

// Acquire blocks until a slot is available or ctx is done.
// Returns ctx.Err() if the context is cancelled before a slot is free.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.ch <- struct{}{}:
		l.acquired.Add(1)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
// Calling Release without a prior successful Acquire is a no-op.
func (l *Limiter) Release() {
	select {
	case <-l.ch:
		l.acquired.Add(-1)
	default:
	}
}

// InFlight returns the number of currently acquired slots.
func (l *Limiter) InFlight() int64 {
	return l.acquired.Load()
}

// Cap returns the maximum concurrency capacity.
func (l *Limiter) Cap() int {
	return cap(l.ch)
}

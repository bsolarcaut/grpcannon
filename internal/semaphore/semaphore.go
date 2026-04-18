// Package semaphore provides a simple counting semaphore backed by a buffered channel.
package semaphore

import (
	"context"
	"fmt"
)

// Semaphore limits concurrent access to a resource.
type Semaphore struct {
	ch chan struct{}
}

// New creates a Semaphore with the given capacity. Returns an error if n < 1.
func New(n int) (*Semaphore, error) {
	if n < 1 {
		return nil, fmt.Errorf("semaphore: capacity must be at least 1, got %d", n)
	}
	return &Semaphore{ch: make(chan struct{}, n)}, nil
}

// Acquire acquires one slot, blocking until one is available or ctx is done.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases one slot back to the semaphore.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
	}
}

// Cap returns the total capacity of the semaphore.
func (s *Semaphore) Cap() int {
	return cap(s.ch)
}

// InUse returns the number of currently acquired slots.
func (s *Semaphore) InUse() int {
	return len(s.ch)
}

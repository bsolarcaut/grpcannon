// Package slot provides a fixed-capacity token bucket that workers
// can acquire before dispatching a request and release afterward.
// Unlike a semaphore, Slot tracks per-slot identifiers so callers
// can correlate in-flight work to a numbered lane.
package slot

import (
	"context"
	"fmt"
	"sync"
)

// Slot represents a numbered lane that has been acquired.
type Slot int

// Pool manages a fixed set of numbered slots.
type Pool struct {
	mu   sync.Mutex
	free []Slot
	ch   chan Slot
}

// New creates a Pool with cap numbered slots (1-indexed).
// cap is clamped to a minimum of 1.
func New(cap int) *Pool {
	if cap < 1 {
		cap = 1
	}
	ch := make(chan Slot, cap)
	for i := 1; i <= cap; i++ {
		ch <- Slot(i)
	}
	return &Pool{ch: ch}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns the acquired Slot or an error.
func (p *Pool) Acquire(ctx context.Context) (Slot, error) {
	select {
	case s := <-p.ch:
		return s, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// Release returns a previously acquired slot back to the pool.
// Releasing an invalid slot (<=0) is a no-op.
func (p *Pool) Release(s Slot) error {
	if s <= 0 {
		return fmt.Errorf("slot: invalid slot %d", s)
	}
	select {
	case p.ch <- s:
		return nil
	default:
		return fmt.Errorf("slot: pool full, slot %d may have been released twice", s)
	}
}

// Cap returns the total capacity of the pool.
func (p *Pool) Cap() int {
	return cap(p.ch)
}

// Available returns the number of slots currently free.
func (p *Pool) Available() int {
	return len(p.ch)
}

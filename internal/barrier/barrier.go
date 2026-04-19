// Package barrier provides a reusable synchronisation barrier that blocks
// goroutines until a target count of arrivals is reached.
package barrier

import (
	"context"
	"sync"
)

// Barrier blocks callers at Wait until exactly n goroutines have arrived.
type Barrier struct {
	mu      sync.Mutex
	cond    *sync.Cond
	target  int
	arrived int
	gen     int // generation counter — resets the barrier for reuse
}

// New returns a Barrier that releases when n goroutines call Wait.
// Panics if n < 1.
func New(n int) *Barrier {
	if n < 1 {
		panic("barrier: n must be >= 1")
	}
	b := &Barrier{target: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks until all n goroutines have called Wait, then releases all of
// them. Returns ctx.Err() if the context is cancelled before the barrier
// opens. On context cancellation the internal counter is decremented so the
// barrier remains usable.
func (b *Barrier) Wait(ctx context.Context) error {
	b.mu.Lock()
	gen := b.gen
	b.arrived++
	if b.arrived == b.target {
		b.gen++
		b.arrived = 0
		b.cond.Broadcast()
		b.mu.Unlock()
		return nil
	}

	// Watch for context cancellation in a separate goroutine.
	doneCh := make(chan struct{})
	defer close(doneCh)
	go func() {
		select {
		case <-ctx.Done():
			b.mu.Lock()
			b.cond.Broadcast()
			b.mu.Unlock()
		case <-doneCh:
		}
	}()

	for b.gen == gen {
		b.cond.Wait()
		if ctx.Err() != nil {
			b.arrived--
			b.mu.Unlock()
			return ctx.Err()
		}
	}
	b.mu.Unlock()
	return nil
}

// Reset resets the barrier to its initial state. Callers currently blocked in
// Wait are released with a nil error.
func (b *Barrier) Reset() {
	b.mu.Lock()
	b.gen++
	b.arrived = 0
	b.cond.Broadcast()
	b.mu.Unlock()
}

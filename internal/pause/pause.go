// Package pause provides a simple run-pause-resume controller that can
// temporarily halt worker goroutines without stopping them entirely.
package pause

import (
	"context"
	"sync"
)

// Controller allows callers to pause and resume a set of workers.
type Controller struct {
	mu      sync.Mutex
	paused  bool
	cond    *sync.Cond
	stopped bool
}

// New returns a ready-to-use Controller.
func New() *Controller {
	c := &Controller{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

// Pause causes subsequent calls to Wait to block.
func (c *Controller) Pause() {
	c.mu.Lock()
	c.paused = true
	c.mu.Unlock()
}

// Resume unblocks all goroutines waiting in Wait.
func (c *Controller) Resume() {
	c.mu.Lock()
	c.paused = false
	c.cond.Broadcast()
	c.mu.Unlock()
}

// Stop permanently unblocks all waiters and prevents future blocking.
func (c *Controller) Stop() {
	c.mu.Lock()
	c.stopped = true
	c.paused = false
	c.cond.Broadcast()
	c.mu.Unlock()
}

// Wait blocks while the controller is paused. It returns ctx.Err() if the
// context is cancelled before the controller is resumed.
func (c *Controller) Wait(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for c.paused && !c.stopped {
		// Use a channel to integrate cond with context cancellation.
		waited := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				c.mu.Lock()
				c.cond.Broadcast()
				c.mu.Unlock()
			case <-waited:
			}
		}()
		c.cond.Wait()
		close(waited)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

// IsPaused reports whether the controller is currently paused.
func (c *Controller) IsPaused() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.paused
}

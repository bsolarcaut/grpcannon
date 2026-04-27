// Package backpressure implements a simple back-pressure controller that
// signals callers to slow down when in-flight work exceeds a high-water mark
// and resumes normal operation once it drops below a low-water mark.
package backpressure

import (
	"context"
	"sync"
)

// Controller tracks in-flight requests and exposes a Wait method that blocks
// new work when the high-water mark is reached.
type Controller struct {
	mu       sync.Mutex
	high     int
	low      int
	inFlight int
	pressure bool
	resume   chan struct{}
}

// New returns a Controller that applies back-pressure above high and releases
// it once in-flight count falls to or below low.
// Panics if high < 1 or low >= high.
func New(high, low int) *Controller {
	if high < 1 {
		panic("backpressure: high must be >= 1")
	}
	if low >= high {
		panic("backpressure: low must be < high")
	}
	if low < 0 {
		low = 0
	}
	return &Controller{
		high:   high,
		low:    low,
		resume: make(chan struct{}),
	}
}

// Acquire increments the in-flight counter. It blocks when back-pressure is
// active until the load drops or ctx is cancelled.
func (c *Controller) Acquire(ctx context.Context) error {
	for {
		c.mu.Lock()
		if !c.pressure {
			c.inFlight++
			c.mu.Unlock()
			return nil
		}
		ch := c.resume
		c.mu.Unlock()

		select {
		case <-ch:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Release decrements the in-flight counter and lifts back-pressure if the
// count has fallen to the low-water mark.
func (c *Controller) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.inFlight > 0 {
		c.inFlight--
	}
	if c.pressure && c.inFlight <= c.low {
		c.pressure = false
		close(c.resume)
		c.resume = make(chan struct{})
	}
}

// Record updates the in-flight gauge directly (useful when Acquire/Release are
// not used) and applies or lifts pressure accordingly.
func (c *Controller) Record(inFlight int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.inFlight = inFlight
	if !c.pressure && c.inFlight >= c.high {
		c.pressure = true
	}
	if c.pressure && c.inFlight <= c.low {
		c.pressure = false
		close(c.resume)
		c.resume = make(chan struct{})
	}
}

// UnderPressure reports whether back-pressure is currently active.
func (c *Controller) UnderPressure() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pressure
}

// InFlight returns the current in-flight count.
func (c *Controller) InFlight() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inFlight
}

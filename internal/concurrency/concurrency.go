// Package concurrency provides dynamic concurrency adjustment
// based on observed error rates and latency signals.
package concurrency

import (
	"sync"
	"time"
)

// Level holds a concurrency reading at a point in time.
type Level struct {
	Workers   int
	RecordedAt time.Time
}

// Controller manages the current concurrency level within bounds.
type Controller struct {
	mu      sync.Mutex
	current int
	min     int
	max     int
	step    int
	history []Level
}

// New creates a Controller with the given initial, min, max, and step values.
// Values are clamped so that min >= 1, initial is within [min, max].
func New(initial, min, max, step int) *Controller {
	if min < 1 {
		min = 1
	}
	if max < min {
		max = min
	}
	if initial < min {
		initial = min
	}
	if initial > max {
		initial = max
	}
	if step < 1 {
		step = 1
	}
	return &Controller{
		current: initial,
		min:     min,
		max:     max,
		step:    step,
	}
}

// Current returns the current concurrency level.
func (c *Controller) Current() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}

// ScaleUp increases concurrency by step, capped at max.
func (c *Controller) ScaleUp() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current += c.step
	if c.current > c.max {
		c.current = c.max
	}
	c.history = append(c.history, Level{Workers: c.current, RecordedAt: time.Now()})
	return c.current
}

// ScaleDown decreases concurrency by step, floored at min.
func (c *Controller) ScaleDown() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current -= c.step
	if c.current < c.min {
		c.current = c.min
	}
	c.history = append(c.history, Level{Workers: c.current, RecordedAt: time.Now()})
	return c.current
}

// History returns a snapshot of recorded concurrency levels.
func (c *Controller) History() []Level {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Level, len(c.history))
	copy(out, c.history)
	return out
}

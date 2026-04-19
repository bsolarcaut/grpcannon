// Package adaptive provides a concurrency controller that adjusts worker
// count up or down based on observed error rate and latency.
package adaptive

import (
	"sync"
	"time"
)

const (
	defaultMin      = 1
	defaultMax      = 512
	defaultStep     = 2
	defaultCooldown = 2 * time.Second
)

// Controller adjusts concurrency dynamically.
type Controller struct {
	mu       sync.Mutex
	current  int
	min      int
	max      int
	step     int
	cooldown time.Duration
	lastAdj  time.Time
}

// Config holds tuning parameters for the controller.
type Config struct {
	Initial  int
	Min      int
	Max      int
	Step     int
	Cooldown time.Duration
}

// New returns a Controller with the given config, applying defaults for zero values.
func New(cfg Config) *Controller {
	if cfg.Min <= 0 {
		cfg.Min = defaultMin
	}
	if cfg.Max <= 0 {
		cfg.Max = defaultMax
	}
	if cfg.Step <= 0 {
		cfg.Step = defaultStep
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = defaultCooldown
	}
	if cfg.Initial <= 0 {
		cfg.Initial = cfg.Min
	}
	return &Controller{
		current:  cfg.Initial,
		min:      cfg.Min,
		max:      cfg.Max,
		step:     cfg.Step,
		cooldown: cfg.Cooldown,
	}
}

// Current returns the current concurrency level.
func (c *Controller) Current() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}

// ScaleUp increases concurrency by one step if cooldown has elapsed.
// Returns the new level.
func (c *Controller) ScaleUp() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Since(c.lastAdj) < c.cooldown {
		return c.current
	}
	c.current += c.step
	if c.current > c.max {
		c.current = c.max
	}
	c.lastAdj = time.Now()
	return c.current
}

// ScaleDown decreases concurrency by one step if cooldown has elapsed.
// Returns the new level.
func (c *Controller) ScaleDown() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Since(c.lastAdj) < c.cooldown {
		return c.current
	}
	c.current -= c.step
	if c.current < c.min {
		c.current = c.min
	}
	c.lastAdj = time.Now()
	return c.current
}

// Package counter provides a thread-safe atomic request counter
// for tracking in-flight and completed requests during a load test run.
package counter

import "sync/atomic"

// Counter holds atomic counters for load test metrics.
type Counter struct {
	total    atomic.Int64
	success  atomic.Int64
	failures atomic.Int64
	inFlight atomic.Int64
}

// New returns a zero-valued Counter.
func New() *Counter {
	return &Counter{}
}

// IncTotal increments the total request count.
func (c *Counter) IncTotal() { c.total.Add(1) }

// IncSuccess increments the success count.
func (c *Counter) IncSuccess() { c.success.Add(1) }

// IncFailure increments the failure count.
func (c *Counter) IncFailure() { c.failures.Add(1) }

// IncInFlight increments the in-flight count.
func (c *Counter) IncInFlight() { c.inFlight.Add(1) }

// DecInFlight decrements the in-flight count.
func (c *Counter) DecInFlight() { c.inFlight.Add(-1) }

// Total returns the total number of requests issued.
func (c *Counter) Total() int64 { return c.total.Load() }

// Success returns the number of successful requests.
func (c *Counter) Success() int64 { return c.success.Load() }

// Failures returns the number of failed requests.
func (c *Counter) Failures() int64 { return c.failures.Load() }

// InFlight returns the current number of in-flight requests.
func (c *Counter) InFlight() int64 { return c.inFlight.Load() }

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.total.Store(0)
	c.success.Store(0)
	c.failures.Store(0)
	c.inFlight.Store(0)
}

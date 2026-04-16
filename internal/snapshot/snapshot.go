// Package snapshot captures a point-in-time view of runner metrics.
package snapshot

import (
	"sync"
	"time"
)

// Metrics holds a moment-in-time view of load-test progress.
type Metrics struct {
	Timestamp  time.Time
	Total      int64
	Successes  int64
	Errors     int64
	InFlight   int64
	Elapsed    time.Duration
}

// Collector accumulates metrics and allows periodic snapshots.
type Collector struct {
	mu        sync.Mutex
	start     time.Time
	total     int64
	successes int64
	errors    int64
	inFlight  int64
}

// NewCollector creates a Collector with the clock started.
func NewCollector() *Collector {
	return &Collector{start: time.Now()}
}

// RecordSuccess increments success and total counters.
func (c *Collector) RecordSuccess() {
	c.mu.Lock()
	c.total++
	c.successes++
	c.mu.Unlock()
}

// RecordError increments error and total counters.
func (c *Collector) RecordError() {
	c.mu.Lock()
	c.total++
	c.errors++
	c.mu.Unlock()
}

// SetInFlight sets the current number of in-flight requests.
func (c *Collector) SetInFlight(n int64) {
	c.mu.Lock()
	c.inFlight = n
	c.mu.Unlock()
}

// Snap returns the current Metrics snapshot.
func (c *Collector) Snap() Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Metrics{
		Timestamp: time.Now(),
		Total:     c.total,
		Successes: c.successes,
		Errors:    c.errors,
		InFlight:  c.inFlight,
		Elapsed:   time.Since(c.start),
	}
}

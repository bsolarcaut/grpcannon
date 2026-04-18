// Package pressure tracks and exposes an instantaneous load score
// in the range [0.0, 1.0] based on in-flight requests vs capacity.
package pressure

import (
	"sync/atomic"
)

// Gauge measures current load pressure.
type Gauge struct {
	capacity int64
	inFlight atomic.Int64
}

// New returns a Gauge with the given capacity.
// Capacity must be >= 1; values below 1 are clamped to 1.
func New(capacity int) *Gauge {
	if capacity < 1 {
		capacity = 1
	}
	return &Gauge{capacity: int64(capacity)}
}

// Acquire records one additional in-flight request.
// It returns false if the gauge is already at or above capacity.
func (g *Gauge) Acquire() bool {
	for {
		cur := g.inFlight.Load()
		if cur >= g.capacity {
			return false
		}
		if g.inFlight.CompareAndSwap(cur, cur+1) {
			return true
		}
	}
}

// Release decrements the in-flight counter.
// Calls below zero are ignored.
func (g *Gauge) Release() {
	for {
		cur := g.inFlight.Load()
		if cur <= 0 {
			return
		}
		if g.inFlight.CompareAndSwap(cur, cur-1) {
			return
		}
	}
}

// Score returns the current pressure as a value in [0.0, 1.0].
func (g *Gauge) Score() float64 {
	cur := g.inFlight.Load()
	if cur <= 0 {
		return 0.0
	}
	if cur >= g.capacity {
		return 1.0
	}
	return float64(cur) / float64(g.capacity)
}

// InFlight returns the current number of in-flight requests.
func (g *Gauge) InFlight() int64 {
	return g.inFlight.Load()
}

// Capacity returns the configured capacity.
func (g *Gauge) Capacity() int64 {
	return g.capacity
}

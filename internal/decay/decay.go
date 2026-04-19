// Package decay provides an exponentially weighted moving average (EWMA)
// suitable for smoothing rate and latency signals over time.
package decay

import (
	"math"
	"sync"
	"time"
)

// EWMA is a thread-safe exponentially weighted moving average.
type EWMA struct {
	mu     sync.Mutex
	alpha  float64
	value  float64
	seeded bool
}

// New returns an EWMA with the given half-life duration.
// A shorter half-life reacts faster to recent changes.
func New(halfLife time.Duration) *EWMA {
	if halfLife <= 0 {
		halfLife = time.Second
	}
	alpha := 1 - math.Exp(-math.Ln2/float64(halfLife))
	return &EWMA{alpha: alpha}
}

// Add incorporates a new sample into the moving average.
func (e *EWMA) Add(v float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.seeded {
		e.value = v
		e.seeded = true
		return
	}
	e.value = e.alpha*v + (1-e.alpha)*e.value
}

// Value returns the current smoothed average.
func (e *EWMA) Value() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}

// Reset clears the EWMA back to an unseeded state.
func (e *EWMA) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.value = 0
	e.seeded = false
}

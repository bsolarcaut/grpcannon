// Package shed implements load shedding based on a configurable error-rate
// threshold. When the observed error rate exceeds the threshold the Shed method
// returns false, signalling callers to drop the request rather than forward it.
package shed

import (
	"sync"
	"time"
)

// Shed decides whether a new request should be accepted.
type Shed struct {
	mu        sync.Mutex
	threshold float64 // 0–1 fraction of errors that triggers shedding
	window    []bool  // circular buffer of recent outcomes
	pos       int
	size      int
	total     int
	errors    int
	cooldown  time.Duration
	lastShed  time.Time
}

// New returns a Shed that starts shedding once errorRate exceeds threshold.
// windowSize is the number of recent calls tracked; cooldown is the minimum
// time between consecutive shed decisions.
func New(threshold float64, windowSize int, cooldown time.Duration) *Shed {
	if windowSize < 1 {
		windowSize = 100
	}
	if cooldown <= 0 {
		cooldown = 100 * time.Millisecond
	}
	return &Shed{
		threshold: threshold,
		window:    make([]bool, windowSize),
		size:      windowSize,
		cooldown:  cooldown,
	}
}

// Record registers the outcome of a completed call (true = error).
func (s *Shed) Record(isError bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old := s.window[s.pos]
	s.window[s.pos] = isError
	s.pos = (s.pos + 1) % s.size

	if s.total < s.size {
		s.total++
	} else if old {
		s.errors--
	}
	if isError {
		s.errors++
	}
}

// Allow returns true when the request should proceed, false when it should be
// shed. It respects the cooldown so shedding decisions are not made too
// frequently.
func (s *Shed) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.total == 0 {
		return true
	}
	rate := float64(s.errors) / float64(s.total)
	if rate < s.threshold {
		return true
	}
	if time.Since(s.lastShed) < s.cooldown {
		return false
	}
	s.lastShed = time.Now()
	return false
}

// Rate returns the current observed error rate (0–1).
func (s *Shed) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.total == 0 {
		return 0
	}
	return float64(s.errors) / float64(s.total)
}

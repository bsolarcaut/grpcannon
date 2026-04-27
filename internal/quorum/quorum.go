// Package quorum tracks whether a minimum fraction of workers are healthy
// before allowing the load test to proceed. If fewer than the required
// fraction of workers report success, Quorum returns an error.
package quorum

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// ErrQuorumNotMet is returned when healthy workers fall below the threshold.
var ErrQuorumNotMet = errors.New("quorum: healthy worker fraction below threshold")

// Quorum tracks worker health votes and evaluates whether the minimum
// healthy fraction has been met.
type Quorum struct {
	threshold float64 // fraction in (0, 1]
	total     atomic.Int64
	healthy   atomic.Int64
}

// New creates a Quorum with the given threshold fraction.
// threshold must be in the range (0, 1]; values outside that range are clamped.
func New(threshold float64) *Quorum {
	if threshold <= 0 {
		threshold = 0.5
	}
	if threshold > 1 {
		threshold = 1
	}
	return &Quorum{threshold: threshold}
}

// Vote records one worker result. healthy=true means the worker is healthy.
func (q *Quorum) Vote(healthy bool) {
	q.total.Add(1)
	if healthy {
		q.healthy.Add(1)
	}
}

// Check returns nil when the healthy fraction meets the threshold, or
// ErrQuorumNotMet otherwise. If no votes have been cast it always returns nil.
func (q *Quorum) Check() error {
	t := q.total.Load()
	if t == 0 {
		return nil
	}
	h := q.healthy.Load()
	fraction := float64(h) / float64(t)
	if fraction < q.threshold {
		return fmt.Errorf("%w: %.2f%% healthy (need %.2f%%)",
			ErrQuorumNotMet, fraction*100, q.threshold*100)
	}
	return nil
}

// Reset clears all recorded votes.
func (q *Quorum) Reset() {
	q.total.Store(0)
	q.healthy.Store(0)
}

// Fraction returns the current healthy fraction. Returns 0 when no votes exist.
func (q *Quorum) Fraction() float64 {
	t := q.total.Load()
	if t == 0 {
		return 0
	}
	return float64(q.healthy.Load()) / float64(t)
}

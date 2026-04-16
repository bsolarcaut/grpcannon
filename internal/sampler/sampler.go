// Package sampler provides probabilistic request sampling for load tests.
package sampler

import (
	"math/rand"
	"sync/atomic"
)

// Sampler decides whether a given request should be recorded in detail.
type Sampler struct {
	rate    float64 // 0.0 – 1.0
	total   atomic.Int64
	sampled atomic.Int64
	rng     *rand.Rand
}

// New returns a Sampler that records approximately rate*100 % of requests.
// A rate of 0 disables sampling; a rate >= 1 samples everything.
func New(rate float64, seed int64) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)), //nolint:gosec
	}
}

// Sample returns true if this request should be sampled.
func (s *Sampler) Sample() bool {
	s.total.Add(1)
	if s.rate == 0 {
		return false
	}
	if s.rate >= 1 {
		s.sampled.Add(1)
		return true
	}
	if s.rng.Float64() < s.rate {
		s.sampled.Add(1)
		return true
	}
	return false
}

// Stats returns total requests seen and how many were sampled.
func (s *Sampler) Stats() (total, sampled int64) {
	return s.total.Load(), s.sampled.Load()
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 { return s.rate }

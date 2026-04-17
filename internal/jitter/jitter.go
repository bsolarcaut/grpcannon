// Package jitter provides utilities for adding randomised jitter to durations,
// which helps spread retry and back-off storms across concurrent workers.
package jitter

import (
	"math/rand"
	"time"
)

// Jitter holds the configuration for jitter calculation.
type Jitter struct {
	rng *rand.Rand
}

// New returns a Jitter seeded with the provided seed.
// Pass a time-based seed in production; a fixed seed in tests.
func New(seed int64) *Jitter {
	//nolint:gosec // weak RNG is acceptable for jitter
	return &Jitter{rng: rand.New(rand.NewSource(seed))}
}

// Full returns a random duration in [0, d).
// If d <= 0 the zero duration is returned.
func (j *Jitter) Full(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(j.rng.Int63n(int64(d)))
}

// Equal returns a random duration in [d/2, d).
// This keeps the delay within half of the base value, avoiding very short waits.
func (j *Jitter) Equal(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	half := d / 2
	return half + time.Duration(j.rng.Int63n(int64(d-half)))
}

// Apply adds jitter to base using the Equal strategy and returns the result
// clamped to [min, max].
func (j *Jitter) Apply(base, min, max time.Duration) time.Duration {
	v := j.Equal(base)
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

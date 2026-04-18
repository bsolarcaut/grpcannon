// Package budget implements an error-budget tracker for load tests.
// It tracks the ratio of failures to total requests and trips when
// the configured threshold is exceeded.
package budget

import (
	"errors"
	"sync/atomic"
)

// ErrExhausted is returned when the error budget is exhausted.
var ErrExhausted = errors.New("error budget exhausted")

// Budget tracks failures against a maximum allowed error rate.
type Budget struct {
	total    atomic.Int64
	failures atomic.Int64
	threshold float64 // 0–1, e.g. 0.05 for 5 %
	minTotal  int64   // minimum requests before tripping
}

// New creates a Budget. threshold is the max failure ratio (0–1).
// minTotal is the minimum number of requests before the budget can trip.
func New(threshold float64, minTotal int64) *Budget {
	if threshold <= 0 {
		threshold = 0.01
	}
	if threshold > 1 {
		threshold = 1
	}
	if minTotal <= 0 {
		minTotal = 10
	}
	return &Budget{threshold: threshold, minTotal: minTotal}
}

// Record registers the outcome of a single request.
func (b *Budget) Record(err error) {
	b.total.Add(1)
	if err != nil {
		b.failures.Add(1)
	}
}

// Check returns ErrExhausted when the error rate exceeds the threshold.
func (b *Budget) Check() error {
	t := b.total.Load()
	if t < b.minTotal {
		return nil
	}
	rate := float64(b.failures.Load()) / float64(t)
	if rate > b.threshold {
		return ErrExhausted
	}
	return nil
}

// Rate returns the current failure ratio.
func (b *Budget) Rate() float64 {
	t := b.total.Load()
	if t == 0 {
		return 0
	}
	return float64(b.failures.Load()) / float64(t)
}

// Snapshot returns total and failure counts.
func (b *Budget) Snapshot() (total, failures int64) {
	return b.total.Load(), b.failures.Load()
}

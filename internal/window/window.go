// Package window provides a sliding-window counter for tracking
// request rates over a rolling time interval.
package window

import (
	"sync"
	"time"
)

// Window is a thread-safe sliding-window counter.
type Window struct {
	mu       sync.Mutex
	buckets  []int64
	size     int
	interval time.Duration
	last     time.Time
}

// New creates a Window with the given number of buckets and per-bucket interval.
// Total window duration = size * interval.
func New(size int, interval time.Duration) *Window {
	if size < 1 {
		size = 1
	}
	return &Window{
		buckets:  make([]int64, size),
		size:     size,
		interval: interval,
		last:     time.Now(),
	}
}

// Add increments the current bucket by delta.
func (w *Window) Add(delta int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	w.buckets[0] += delta
}

// Sum returns the total count across all buckets.
func (w *Window) Sum() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	var total int64
	for _, v := range w.buckets {
		total += v
	}
	return total
}

// Rate returns requests per second averaged over the full window.
func (w *Window) Rate() float64 {
	total := w.Sum()
	secs := w.interval.Seconds() * float64(w.size)
	if secs == 0 {
		return 0
	}
	return float64(total) / secs
}

// advance rotates buckets based on elapsed time.
func (w *Window) advance(now time.Time) {
	elapsed := now.Sub(w.last)
	steps := int(elapsed / w.interval)
	if steps <= 0 {
		return
	}
	if steps > w.size {
		steps = w.size
	}
	for i := w.size - 1; i >= steps; i-- {
		w.buckets[i] = w.buckets[i-steps]
	}
	for i := 0; i < steps && i < w.size; i++ {
		w.buckets[i] = 0
	}
	w.last = w.last.Add(time.Duration(steps) * w.interval)
}

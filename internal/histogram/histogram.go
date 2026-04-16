// Package histogram provides a simple latency histogram bucketing utility.
package histogram

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Bucket represents a single latency range and its count.
type Bucket struct {
	Low   time.Duration
	High  time.Duration
	Count int
}

// Histogram holds latency buckets.
type Histogram struct {
	buckets []Bucket
	bounds  []time.Duration
}

// New creates a Histogram with the given boundary durations.
// Each pair of adjacent bounds forms a bucket; an overflow bucket is appended.
func New(bounds []time.Duration) *Histogram {
	buckets := make([]Bucket, len(bounds))
	for i, b := range bounds {
		low := time.Duration(0)
		if i > 0 {
			low = bounds[i-1]
		}
		buckets[i] = Bucket{Low: low, High: b}
	}
	// overflow bucket
	buckets = append(buckets, Bucket{Low: bounds[len(bounds)-1], High: -1})
	return &Histogram{buckets: buckets, bounds: bounds}
}

// Record adds a duration to the appropriate bucket.
func (h *Histogram) Record(d time.Duration) {
	for i, bound := range h.bounds {
		if d < bound {
			h.buckets[i].Count++
			return
		}
	}
	h.buckets[len(h.buckets)-1].Count++
}

// Buckets returns a copy of the internal buckets.
func (h *Histogram) Buckets() []Bucket {
	out := make([]Bucket, len(h.buckets))
	copy(out, h.buckets)
	return out
}

// Print writes an ASCII bar chart to w.
func (h *Histogram) Print(w io.Writer) {
	max := 1
	for _, b := range h.buckets {
		if b.Count > max {
			max = b.Count
		}
	}
	for _, b := range h.buckets {
		label := fmt.Sprintf("%6s - %-6s", b.Low.Round(time.Millisecond), b.High.Round(time.Millisecond))
		if b.High < 0 {
			label = fmt.Sprintf("%6s - %-6s", b.Low.Round(time.Millisecond), "+Inf  ")
		}
		bar := strings.Repeat("█", b.Count*40/max)
		fmt.Fprintf(w, "%s | %-40s %d\n", label, bar, b.Count)
	}
}

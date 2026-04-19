// Package scoreboard tracks per-worker request counts and latencies
// for real-time display during a load test run.
package scoreboard

import (
	"sync"
	"time"
)

// Entry holds aggregated stats for a single worker.
type Entry struct {
	WorkerID int
	Total    int64
	Errors   int64
	TotalLat time.Duration
}

// AvgLatency returns the mean latency for the worker, or zero if no requests.
func (e Entry) AvgLatency() time.Duration {
	if e.Total == 0 {
		return 0
	}
	return e.TotalLat / time.Duration(e.Total)
}

// Board is a concurrency-safe scoreboard for worker stats.
type Board struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an empty Board.
func New() *Board {
	return &Board{entries: make(map[int]*Entry)}
}

// Record adds a single call result for the given worker.
func (b *Board) Record(workerID int, lat time.Duration, err bool) {
	b.mu.Lock()
	e, ok := b.entries[workerID]
	if !ok {
		e = &Entry{WorkerID: workerID}
		b.entries[workerID] = e
	}
	e.Total++
	e.TotalLat += lat
	if err {
		e.Errors++
	}
	b.mu.Unlock()
}

// Snapshot returns a copy of all entries.
func (b *Board) Snapshot() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all recorded data.
func (b *Board) Reset() {
	b.mu.Lock()
	b.entries = make(map[int]*Entry)
	b.mu.Unlock()
}

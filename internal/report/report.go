// Package report assembles a final run report from collected stats and
// snapshot data, ready for rendering by the reporter.
package report

import (
	"time"

	"github.com/example/grpcannon/internal/stats"
)

// RunReport holds everything produced by a completed load-test run.
type RunReport struct {
	Target    string
	Method    string
	StartedAt time.Time
	Duration  time.Duration
	Summary   stats.Summary
	Workers   int
	RPS       float64
}

// Builder accumulates run metadata and produces a RunReport.
type Builder struct {
	target    string
	method    string
	workers   int
	started   time.Time
	results   []stats.Result
}

// New returns a Builder initialised with run metadata.
func New(target, method string, workers int) *Builder {
	return &Builder{
		target:  target,
		method:  method,
		workers: workers,
		started: time.Now(),
	}
}

// Add appends a single call result to the builder.
func (b *Builder) Add(r stats.Result) {
	b.results = append(b.results, r)
}

// Build computes statistics and returns the final RunReport.
func (b *Builder) Build() RunReport {
	duration := time.Since(b.started)
	summary := stats.Compute(b.results)

	rps := 0.0
	if secs := duration.Seconds(); secs > 0 {
		rps = float64(summary.Total) / secs
	}

	return RunReport{
		Target:   b.target,
		Method:   b.method,
		StartedAt: b.started,
		Duration: duration,
		Summary:  summary,
		Workers:  b.workers,
		RPS:      rps,
	}
}

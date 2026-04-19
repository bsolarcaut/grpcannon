// Package observe provides a real-time metrics observer that periodically
// samples a snapshot and emits structured log lines during a load run.
package observe

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nickpoorman/grpcannon/internal/snapshot"
)

// Observer samples a [snapshot.Collector] at a fixed interval and writes
// one-line summaries to an [io.Writer].
type Observer struct {
	collector *snapshot.Collector
	interval  time.Duration
	out       io.Writer
}

// Option configures an Observer.
type Option func(*Observer)

// WithInterval overrides the default 1 s sampling interval.
func WithInterval(d time.Duration) Option {
	return func(o *Observer) {
		if d > 0 {
			o.interval = d
		}
	}
}

// WithWriter redirects output away from os.Stdout.
func WithWriter(w io.Writer) Option {
	return func(o *Observer) {
		if w != nil {
			o.out = w
		}
	}
}

// New creates an Observer backed by the given Collector.
func New(c *snapshot.Collector, opts ...Option) *Observer {
	o := &Observer{
		collector: c,
		interval:  time.Second,
		out:       os.Stdout,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Run starts the sampling loop and blocks until ctx is cancelled.
func (o *Observer) Run(ctx context.Context) {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s := o.collector.Snap()
			f
			=%.1fs totald ok=%d=%d inflight=%d\.Elapsed.Seconds(), s.Total, s.Success, s.Errors, s.InFlight,
			)
		}
	}
}

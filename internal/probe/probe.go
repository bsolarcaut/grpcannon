// Package probe provides a lightweight health-check mechanism that sends
// a single no-payload gRPC request and reports whether the target is ready.
package probe

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	OK       bool
	Latency  time.Duration
	Err      error
}

// Prober checks whether a gRPC target is reachable and responsive.
type Prober struct {
	conn    *grpc.ClientConn
	method  string
	timeout time.Duration
}

const defaultTimeout = 5 * time.Second

// New creates a Prober for the given connection and fully-qualified method.
func New(conn *grpc.ClientConn, method string, timeout time.Duration) (*Prober, error) {
	if conn == nil {
		return nil, fmt.Errorf("probe: conn must not be nil")
	}
	if method == "" {
		return nil, fmt.Errorf("probe: method must not be empty")
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Prober{conn: conn, method: method, timeout: timeout}, nil
}

// Check performs a single probe and returns a Result.
func (p *Prober) Check(ctx context.Context) Result {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	start := time.Now()
	err := p.conn.Invoke(ctx, p.method, struct{}{}, &struct{}{})
	latency := time.Since(start)

	if err != nil {
		return Result{OK: false, Latency: latency, Err: err}
	}
	return Result{OK: true, Latency: latency}
}

// CheckN runs up to n probes, stopping on the first success.
func (p *Prober) CheckN(ctx context.Context, n int, interval time.Duration) Result {
	var last Result
	for i := 0; i < n; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return Result{OK: false, Err: ctx.Err()}
			case <-time.After(interval):
			}
		}
		last = p.Check(ctx)
		if last.OK {
			return last
		}
	}
	return last
}

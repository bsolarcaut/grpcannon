package worker

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// CallFunc is a function that performs a single gRPC call and returns an error.
type CallFunc func(ctx context.Context, conn *grpc.ClientConn) error

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration time.Duration
	Err      error
}

// Worker sends requests using the provided CallFunc and reports results.
type Worker struct {
	conn     *grpc.ClientConn
	callFn   CallFunc
	results  chan<- Result
}

// New creates a new Worker.
func New(conn *grpc.ClientConn, fn CallFunc, results chan<- Result) *Worker {
	return &Worker{
		conn:    conn,
		callFn:  fn,
		results: results,
	}
}

// Run executes requests in a loop until the context is cancelled or n requests
// have been sent. If n <= 0, it runs until context cancellation.
func (w *Worker) Run(ctx context.Context, n int) {
	count := 0
	for {
		if ctx.Err() != nil {
			return
		}
		if n > 0 && count >= n {
			return
		}

		start := time.Now()
		err := w.callFn(ctx, w.conn)
		dur := time.Since(start)

		select {
		case w.results <- Result{Duration: dur, Err: err}:
		case <-ctx.Done():
			return
		}
		count++
	}
}

package runner

import (
	"context"
	"sync"
	"time"

	"github.com/yourusername/grpcannon/internal/config"
)

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration time.Duration
	Err      error
}

// Runner orchestrates concurrent gRPC load testing.
type Runner struct {
	cfg    *config.Config
	caller Caller
}

// Caller is the interface for making a single gRPC request.
type Caller interface {
	Call(ctx context.Context) error
}

// New creates a new Runner with the given config and caller.
func New(cfg *config.Config, caller Caller) *Runner {
	return &Runner{cfg: cfg, caller: caller}
}

// Run executes the load test and returns a slice of Results.
func (r *Runner) Run(ctx context.Context) []Result {
	results := make([]Result, 0, r.cfg.TotalRequests)
	resultCh := make(chan Result, r.cfg.TotalRequests)

	sem := make(chan struct{}, r.cfg.Concurrency)
	var wg sync.WaitGroup

	for i := 0; i < r.cfg.TotalRequests; i++ {
		select {
		case <-ctx.Done():
			goto done
		default:
		}

		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			start := time.Now()
			err := r.caller.Call(ctx)
			resultCh <- Result{
				Duration: time.Since(start),
				Err:      err,
			}
		}()
	}

done:
	wg.Wait()
	close(resultCh)

	for res := range resultCh {
		results = append(results, res)
	}
	return results
}

// Package warmup provides a pre-load warm-up phase that sends a small
// number of requests before the main benchmark run begins.
package warmup

import (
	"context"
	"fmt"
	"time"
)

// Caller is the interface satisfied by invoker.Invoker.
type Caller interface {
	Call(ctx context.Context) error
}

// Config holds warm-up parameters.
type Config struct {
	Requests    int
	Concurrency int
	Timeout     time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Requests:    10,
		Concurrency: 1,
		Timeout:     5 * time.Second,
	}
}

// Run executes the warm-up phase, firing cfg.Requests sequential calls.
// It returns an error if more than half of the calls fail.
func Run(ctx context.Context, caller Caller, cfg Config) error {
	if cfg.Requests <= 0 {
		return nil
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultConfig().Timeout
	}

	failed := 0
	for i := 0; i < cfg.Requests; i++ {
		callCtx, cancel := context.WithTimeout(ctx, timeout)
		err := caller.Call(callCtx)
		cancel()
		if err != nil {
			failed++
		}
	}

	if failed > cfg.Requests/2 {
		return fmt.Errorf("warmup: %d/%d requests failed", failed, cfg.Requests)
	}
	return nil
}

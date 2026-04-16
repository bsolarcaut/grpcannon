// Package retry provides configurable retry logic for gRPC call failures.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines how retries are performed.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	RetryOn     func(err error) bool
}

// DefaultPolicy returns a Policy that retries up to 3 times with a 100ms delay
// on any non-nil, non-context error.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
		RetryOn:     IsRetryable,
	}
}

// IsRetryable returns true if the error is neither a context cancellation
// nor a context deadline exceeded error.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

// Do executes fn according to p, retrying on eligible errors.
// It returns the last error if all attempts are exhausted.
func (p Policy) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	var err error
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn(ctx)
		if err == nil {
			return nil
		}
		retryOn := p.RetryOn
		if retryOn == nil {
			retryOn = IsRetryable
		}
		if !retryOn(err) {
			return err
		}
		if attempt < p.MaxAttempts-1 && p.Delay > 0 {
			select {
			case <-time.After(p.Delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return err
}

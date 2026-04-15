// Package backoff provides simple retry backoff strategies for gRPC call retries.
package backoff

import (
	"context"
	"math"
	"time"
)

// Strategy defines how long to wait before the next retry attempt.
type Strategy interface {
	Wait(attempt int) time.Duration
}

// Exponential implements an exponential backoff strategy with optional jitter.
type Exponential struct {
	// Base is the initial wait duration.
	Base time.Duration
	// Max caps the computed wait duration.
	Max time.Duration
	// Multiplier is the growth factor applied per attempt (default 2.0).
	Multiplier float64
}

// DefaultExponential returns an Exponential backoff with sensible defaults.
func DefaultExponential() *Exponential {
	return &Exponential{
		Base:       50 * time.Millisecond,
		Max:        5 * time.Second,
		Multiplier: 2.0,
	}
}

// Wait returns the duration to sleep before attempt n (0-indexed).
func (e *Exponential) Wait(attempt int) time.Duration {
	mul := e.Multiplier
	if mul <= 0 {
		mul = 2.0
	}
	d := float64(e.Base) * math.Pow(mul, float64(attempt))
	if d > float64(e.Max) {
		d = float64(e.Max)
	}
	return time.Duration(d)
}

// Do executes fn up to maxAttempts times, sleeping according to s between
// retries. It stops early if ctx is cancelled and returns ctx.Err().
func Do(ctx context.Context, s Strategy, maxAttempts int, fn func() error) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn()
		if err == nil {
			return nil
		}
		if i == maxAttempts-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.Wait(i)):
		}
	}
	return err
}

// Package ratelimiter provides a token-bucket rate limiter for controlling
// the request rate during load testing.
package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter controls the rate at which requests are dispatched.
type RateLimiter struct {
	ticker *time.Ticker
	done   chan struct{}
}

// New creates a RateLimiter that allows up to rps requests per second.
// If rps is zero or negative, no rate limiting is applied (returns nil).
func New(rps int) *RateLimiter {
	if rps <= 0 {
		return nil
	}
	interval := time.Second / time.Duration(rps)
	return &RateLimiter{
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}
}

// Wait blocks until the rate limiter permits the next request or the context
// is cancelled. Returns ctx.Err() if the context is done before a token is
// available, or nil on success.
func (r *RateLimiter) Wait(ctx context.Context) error {
	if r == nil {
		return nil
	}
	select {
	case <-r.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases the resources held by the RateLimiter.
func (r *RateLimiter) Stop() {
	if r == nil {
		return
	}
	r.ticker.Stop()
}

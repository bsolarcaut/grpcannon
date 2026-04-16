// Package timeout provides helpers for deriving per-call deadlines.
package timeout

import (
	"context"
	"errors"
	"time"
)

// ErrZeroDuration is returned when a zero or negative duration is provided.
var ErrZeroDuration = errors.New("timeout: duration must be positive")

// WithDeadline returns a child context that is cancelled after d, along with
// its cancel function. The caller must call cancel to release resources.
func WithDeadline(parent context.Context, d time.Duration) (context.Context, context.CancelFunc, error) {
	if d <= 0 {
		return nil, nil, ErrZeroDuration
	}
	ctx, cancel := context.WithTimeout(parent, d)
	return ctx, cancel, nil
}

// Clamp returns d clamped to [min, max]. If min > max it returns min.
func Clamp(d, min, max time.Duration) time.Duration {
	if d < min {
		return min
	}
	if max > 0 && d > max {
		return max
	}
	return d
}

// Default returns d if d > 0, otherwise returns fallback.
func Default(d, fallback time.Duration) time.Duration {
	if d > 0 {
		return d
	}
	return fallback
}

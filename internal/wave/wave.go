// Package wave provides a stepped load-wave scheduler that ramps concurrency
// up to a peak, holds it for a soak duration, then ramps back down. It is
// useful for stress-testing a service with a controlled burst profile.
package wave

import (
	"context"
	"time"
)

// Config holds the parameters for a single wave.
type Config struct {
	// Start is the initial worker count at the beginning of the ramp-up.
	Start int
	// Peak is the maximum worker count reached at the top of the wave.
	Peak int
	// Steps is the number of increments used to climb from Start to Peak
	// (and to descend symmetrically). Must be >= 1.
	Steps int
	// StepDuration is how long the scheduler waits at each step level.
	StepDuration time.Duration
	// SoakDuration is how long the scheduler holds at Peak before descending.
	SoakDuration time.Duration
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Start:        1,
		Peak:         10,
		Steps:        5,
		StepDuration: 2 * time.Second,
		SoakDuration: 10 * time.Second,
	}
}

// Schedule emits the target worker count at each step of the wave onto the
// returned channel. The channel is closed when the wave completes or ctx is
// cancelled. The caller is responsible for adjusting the actual concurrency
// level in response to each emitted value.
//
// Sequence: ramp-up steps → soak → ramp-down steps → channel close.
func Schedule(ctx context.Context, cfg Config) <-chan int {
	ch := make(chan int, 1)

	go func() {
		defer close(ch)

		steps := cfg.Steps
		if steps < 1 {
			steps = 1
		}
		start := cfg.Start
		if start < 1 {
			start = 1
		}
		peak := cfg.Peak
		if peak < start {
			peak = start
		}

		spread := peak - start

		// emit is a helper that sends v or returns false if ctx is done.
		emit := func(v int) bool {
			select {
			case ch <- v:
				return true
			case <-ctx.Done():
				return false
			}
		}

		// sleep is a helper that waits d or returns false if ctx is done.
		sleep := func(d time.Duration) bool {
			if d <= 0 {
				return true
			}
			select {
			case <-time.After(d):
				return true
			case <-ctx.Done():
				return false
			}
		}

		// Ramp up.
		for i := 0; i <= steps; i++ {
			level := start + (spread*i)/steps
			if !emit(level) {
				return
			}
			if i < steps {
				if !sleep(cfg.StepDuration) {
					return
				}
			}
		}

		// Soak at peak.
		if !sleep(cfg.SoakDuration) {
			return
		}

		// Ramp down.
		for i := steps - 1; i >= 0; i-- {
			level := start + (spread*i)/steps
			if !emit(level) {
				return
			}
			if i > 0 {
				if !sleep(cfg.StepDuration) {
					return
				}
			}
		}
	}()

	return ch
}

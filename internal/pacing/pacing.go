// Package pacing provides a closed-loop request pacing controller that
// adjusts the inter-request delay to target a desired requests-per-second
// rate while smoothing out bursts using an exponentially weighted moving
// average of the observed throughput.
package pacing

import (
	"context"
	"sync"
	"time"
)

const (
	// defaultAlpha is the EWMA smoothing factor (0 < α ≤ 1).
	// Lower values produce a smoother but slower-reacting estimate.
	defaultAlpha = 0.2

	// minDelay is the floor for the computed inter-request delay.
	minDelay = time.Microsecond
)

// Pacer controls the rate at which callers are allowed to proceed.
// It is safe for concurrent use.
type Pacer struct {
	mu       sync.Mutex
	targetRPS float64
	alpha     float64
	ewmaRPS   float64
	last      time.Time
}

// New creates a Pacer targeting the given requests-per-second rate.
// targetRPS must be positive; values ≤ 0 are clamped to 1.
func New(targetRPS float64) *Pacer {
	if targetRPS <= 0 {
		targetRPS = 1
	}
	return &Pacer{
		targetRPS: targetRPS,
		alpha:     defaultAlpha,
	}
}

// WithAlpha returns a new Pacer with the given EWMA smoothing factor.
// alpha is clamped to the range (0, 1].
func (p *Pacer) WithAlpha(alpha float64) *Pacer {
	if alpha <= 0 {
		alpha = defaultAlpha
	}
	if alpha > 1 {
		alpha = 1
	}
	p.mu.Lock()
	p.alpha = alpha
	p.mu.Unlock()
	return p
}

// Wait blocks until the pacer determines the caller may proceed, or until
// ctx is cancelled. It updates the internal EWMA on each call so subsequent
// delays self-correct toward the target RPS.
func (p *Pacer) Wait(ctx context.Context) error {
	delay := p.next()
	if delay <= 0 {
		return nil
	}
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// next computes the delay before the next request and updates internal state.
func (p *Pacer) next() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.last.IsZero() {
		// First call — seed the EWMA with the target rate and proceed immediately.
		p.last = now
		p.ewmaRPS = p.targetRPS
		return 0
	}

	elapsed := now.Sub(p.last).Seconds()
	p.last = now

	// Compute the instantaneous RPS from the elapsed time since last call.
	var instantRPS float64
	if elapsed > 0 {
		instantRPS = 1.0 / elapsed
	}

	// Update EWMA.
	p.ewmaRPS = p.alpha*instantRPS + (1-p.alpha)*p.ewmaRPS

	// Derive the next delay from the error between target and observed rate.
	// If we are running too fast (ewma > target) we sleep longer; if too slow
	// (ewma < target) we shrink the delay.
	desiredInterval := time.Duration(float64(time.Second) / p.targetRPS)

	if p.ewmaRPS <= 0 {
		return desiredInterval
	}

	// Scale the desired interval by the ratio of observed to target RPS so the
	// pacer self-corrects.
	scale := p.ewmaRPS / p.targetRPS
	adjusted := time.Duration(float64(desiredInterval) * scale)
	if adjusted < minDelay {
		adjusted = minDelay
	}
	return adjusted
}

// RPS returns the current EWMA-smoothed observed requests-per-second.
func (p *Pacer) RPS() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ewmaRPS
}

// SetTarget updates the target RPS at runtime. Values ≤ 0 are ignored.
func (p *Pacer) SetTarget(rps float64) {
	if rps <= 0 {
		return
	}
	p.mu.Lock()
	p.targetRPS = rps
	p.mu.Unlock()
}

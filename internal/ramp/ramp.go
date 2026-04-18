// Package ramp provides a concurrency ramp-up strategy that gradually
// increases the number of active workers over a warm-up duration.
package ramp

import (
	"context"
	"time"
)

// Config holds ramp-up parameters.
type Config struct {
	// StartWorkers is the initial concurrency level.
	StartWorkers int
	// EndWorkers is the target concurrency level.
	EndWorkers int
	// Duration is the total time to ramp from Start to End.
	Duration time.Duration
	// Steps is the number of increments to apply over Duration.
	Steps int
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		StartWorkers: 1,
		EndWorkers:   10,
		Duration:     10 * time.Second,
		Steps:        10,
	}
}

// Schedule returns a channel that emits the concurrency level at each step.
// The channel is closed when the ramp is complete or the context is cancelled.
func Schedule(ctx context.Context, cfg Config) <-chan int {
	ch := make(chan int, cfg.Steps)
	go func() {
		defer close(ch)
		if cfg.Steps <= 0 || cfg.Duration <= 0 {
			select {
			case ch <- cfg.EndWorkers:
			case <-ctx.Done():
			}
			return
		}
		interval := cfg.Duration / time.Duration(cfg.Steps)
		delta := float64(cfg.EndWorkers-cfg.StartWorkers) / float64(cfg.Steps)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for step := 0; step < cfg.Steps; step++ {
			workers := cfg.StartWorkers + int(float64(step)*delta+0.5)
			if step == cfg.Steps-1 {
				workers = cfg.EndWorkers
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case ch <- workerscase <-ctx.Done():
					return
				}
		}
	}n}

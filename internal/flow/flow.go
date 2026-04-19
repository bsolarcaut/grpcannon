// Package flow orchestrates the full load-test lifecycle: warmup,
// ramp-up, steady-state, and graceful drain.
package flow

import (
	"context"
	"fmt"
	"time"

	"github.com/nickpoorman/grpcannon/internal/logger"
	"github.com/nickpoorman/grpcannon/internal/ramp"
	"github.com/nickpoorman/grpcannon/internal/runner"
	"github.com/nickpoorman/grpcannon/internal/warmup"
)

// Config controls the flow execution.
type Config struct {
	Warmup   warmup.Config
	Ramp     ramp.Config
	Steady   time.Duration // how long to hold peak concurrency
	Runner   *runner.Runner
	Log      *logger.Logger
}

// Run executes warmup → ramp → steady → drain in sequence.
// It returns the first non-nil error encountered.
func Run(ctx context.Context, cfg Config) error {
	log := cfg.Log
	if log == nil {
		log = logger.Default()
	}

	log.Info("flow: starting warmup")
	if err := warmup.Run(ctx, cfg.Warmup); err != nil {
		return fmt.Errorf("flow warmup: %w", err)
	}

	log.Info("flow: starting ramp")
	steps := make(chan ramp.Step)
	rampCtx, cancelRamp := context.WithCancel(ctx)
	defer cancelRamp()

	go ramp.Schedule(rampCtx, cfg.Ramp, steps)

	for s := range steps {
		log.Info(fmt.Sprintf("flow: ramp step workers=%d", s.Workers))
		_ = s // concurrency adjustment handled by adaptive controller
	}

	log.Info(fmt.Sprintf("flow: steady state for %s", cfg.Steady))
	steadyCtx, cancelSteady := context.WithTimeout(ctx, cfg.Steady)
	defer cancelSteady()

	if err := cfg.Runner.Run(steadyCtx); err != nil && steadyCtx.Err() == nil {
		return fmt.Errorf("flow steady: %w", err)
	}

	log.Info("flow: complete")
	return nil
}

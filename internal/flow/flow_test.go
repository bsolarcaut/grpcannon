package flow_test

import (
	"context"
	"testing"
	"time"

	"github.com/nickpoorman/grpcannon/internal/flow"
	"github.com/nickpoorman/grpcannon/internal/ramp"
	"github.com/nickpoorman/grpcannon/internal/runner"
	"github.com/nickpoorman/grpcannon/internal/warmup"
	"github.com/nickpoorman/grpcannon/internal/config"
)

func noop(ctx context.Context) error { return nil }

func baseFlowConfig(t *testing.T) flow.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Target = "localhost:50051"
	cfg.Method = "/svc/Method"
	r, err := runner.New(cfg, noop)
	if err != nil {
		t.Fatalf("runner.New: %v", err)
	}
	return flow.Config{
		Warmup: warmup.Config{
			Requests: 0,
			Call:     noop,
		},
		Ramp: ramp.Config{
			Start:    1,
			End:      1,
			Step:     1,
			Interval: time.Millisecond,
		},
		Steady: 10 * time.Millisecond,
		Runner: r,
	}
}

func TestFlow_Run_Completes(t *testing.T) {
	cfg := baseFlowConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := flow.Run(ctx, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFlow_Run_ContextCancellation(t *testing.T) {
	cfg := baseFlowConfig(t)
	cfg.Steady = 10 * time.Second // long steady state
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	// Should not hang
	_ = flow.Run(ctx, cfg)
}

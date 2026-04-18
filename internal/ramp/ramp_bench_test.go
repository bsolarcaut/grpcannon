package ramp_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ramp"
)

func BenchmarkSchedule_Drain(b *testing.B) {
	cfg := ramp.Config{
		StartWorkers: 1,
		EndWorkers:   b.N + 1,
		Duration:     time.Duration(b.N) * time.Microsecond,
		Steps:        b.N,
	}
	if cfg.Steps == 0 {
		cfg.Steps = 1
	}
	b.ResetTimer()
	ctx := context.Background()
	ch := ramp.Schedule(ctx, cfg)
	for range ch {
	}
}

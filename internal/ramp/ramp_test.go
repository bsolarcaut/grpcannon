package ramp_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ramp"
)

func TestDefault_Fields(t *testing.T) {
	cfg := ramp.Default()
	if cfg.StartWorkers <= 0 {
		t.Fatalf("expected positive StartWorkers, got %d", cfg.StartWorkers)
	}
	if cfg.EndWorkers <= cfg.StartWorkers {
		t.Fatalf("expected EndWorkers > StartWorkers")
	}
	if cfg.Steps <= 0 {
		t.Fatalf("expected positive Steps")
	}
}

func TestSchedule_EmitsEndWorkersImmediately(t *testing.T) {
	cfg := ramp.Config{StartWorkers: 1, EndWorkers: 5, Duration: 0, Steps: 0}
	ctx := context.Background()
	ch := ramp.Schedule(ctx, cfg)
	val := <-ch
	if val != 5 {
		t.Fatalf("expected 5, got %d", val)
	}
}

func TestSchedule_CorrectStepCount(t *testing.T) {
	cfg := ramp.Config{
		StartWorkers: 1,
		EndWorkers:   4,
		Duration:     40 * time.Millisecond,
		Steps:        4,
	}
	ctx := context.Background()
	ch := ramp.Schedule(ctx, cfg)
	var levels []int
	for l := range ch {
		levels = append(levels, l)
	}
	if len(levels) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(levels))
	}
	if levels[len(levels)-1] != cfg.EndWorkers {
		t.Fatalf("last step should equal EndWorkers")
	}
}

func TestSchedule_ContextCancellation(t *testing.T) {
	cfg := ramp.Config{
		StartWorkers: 1,
		EndWorkers:   20,
		Duration:     10 * time.Second,
		Steps:        20,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	ch := ramp.Schedule(ctx, cfg)
	var count int
	for range ch {
		count++
	}
	if count >= 20 {
		t.Fatalf("expected cancellation to stop ramp early, got %d steps", count)
	}
}

func TestSchedule_MonotonicallyIncreasing(t *testing.T) {
	cfg := ramp.Config{
		StartWorkers: 2,
		EndWorkers:   6,
		Duration:     40 * time.Millisecond,
		Steps:        4,
	}
	ctx := context.Background()
	ch := ramp.Schedule(ctx, cfg)
	prev := 0
	for l := range ch {
		if l < prev {
			t.Fatalf("non-monotonic step: %d after %d", l, prev)
		}
		prev = l
	}
}

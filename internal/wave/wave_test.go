package wave_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/wave"
)

func TestDefault_Fields(t *testing.T) {
	cfg := wave.Default()

	if cfg.Start <= 0 {
		t.Fatalf("expected Start > 0, got %d", cfg.Start)
	}
	if cfg.Peak <= cfg.Start {
		t.Fatalf("expected Peak > Start, got Peak=%d Start=%d", cfg.Peak, cfg.Start)
	}
	if cfg.Period <= 0 {
		t.Fatalf("expected Period > 0, got %v", cfg.Period)
	}
	if cfg.Step <= 0 {
		t.Fatalf("expected Step > 0, got %d", cfg.Step)
	}
}

func TestSchedule_EmitsInitialWorkers(t *testing.T) {
	cfg := wave.Default()
	cfg.Start = 2
	cfg.Peak = 4
	cfg.Period = 100 * time.Millisecond
	cfg.Step = 2

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := wave.Schedule(ctx, cfg)

	select {
	case v, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before first value")
		}
		if v <= 0 {
			t.Fatalf("expected positive worker count, got %d", v)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for first schedule tick")
	}
}

func TestSchedule_ContextCancellation(t *testing.T) {
	cfg := wave.Default()
	cfg.Period = 10 * time.Second // long period so it doesn't complete naturally

	ctx, cancel := context.WithCancel(context.Background())
	ch := wave.Schedule(ctx, cfg)

	// drain the first value
	select {
	case <-ch:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for initial value")
	}

	cancel()

	// channel should close after cancellation
	select {
	case _, ok := <-ch:
		if ok {
			// drain remaining buffered values
			for range ch {
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel did not close after context cancellation")
	}
}

func TestSchedule_OscillatesBetweenStartAndPeak(t *testing.T) {
	cfg := wave.Default()
	cfg.Start = 2
	cfg.Peak = 6
	cfg.Period = 80 * time.Millisecond
	cfg.Step = 2

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := wave.Schedule(ctx, cfg)

	var values []int
	timeout := time.After(2 * time.Second)
collect:
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				break collect
			}
			values = append(values, v)
			if len(values) >= 8 {
				break collect
			}
		case <-timeout:
			break collect
		}
	}

	if len(values) == 0 {
		t.Fatal("no values emitted")
	}

	for _, v := range values {
		if v < cfg.Start || v > cfg.Peak {
			t.Errorf("value %d out of range [%d, %d]", v, cfg.Start, cfg.Peak)
		}
	}
}

func TestSchedule_ReachesPeak(t *testing.T) {
	cfg := wave.Default()
	cfg.Start = 1
	cfg.Peak = 4
	cfg.Period = 60 * time.Millisecond
	cfg.Step = 1

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := wave.Schedule(ctx, cfg)

	peakSeen := false
	timeout := time.After(2 * time.Second)
collect:
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				break collect
			}
			if v == cfg.Peak {
				peakSeen = true
				break collect
			}
		case <-timeout:
			break collect
		}
	}

	if !peakSeen {
		t.Errorf("peak value %d was never emitted", cfg.Peak)
	}
}

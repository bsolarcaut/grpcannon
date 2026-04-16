package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/lukasgolino/grpcannon/internal/timeout"
)

func TestWithDeadline_ZeroDuration(t *testing.T) {
	_, _, err := timeout.WithDeadline(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestWithDeadline_NegativeDuration(t *testing.T) {
	_, _, err := timeout.WithDeadline(context.Background(), -time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestWithDeadline_ValidDuration(t *testing.T) {
	ctx, cancel, err := timeout.WithDeadline(context.Background(), 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("deadline should be in the future")
	}
}

func TestClamp_BelowMin(t *testing.T) {
	result := timeout.Clamp(10*time.Millisecond, 100*time.Millisecond, time.Second)
	if result != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", result)
	}
}

func TestClamp_AboveMax(t *testing.T) {
	result := timeout.Clamp(5*time.Second, 100*time.Millisecond, time.Second)
	if result != time.Second {
		t.Fatalf("expected 1s, got %v", result)
	}
}

func TestClamp_WithinRange(t *testing.T) {
	result := timeout.Clamp(500*time.Millisecond, 100*time.Millisecond, time.Second)
	if result != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", result)
	}
}

func TestDefault_UsesProvided(t *testing.T) {
	result := timeout.Default(2*time.Second, 5*time.Second)
	if result != 2*time.Second {
		t.Fatalf("expected 2s, got %v", result)
	}
}

func TestDefault_UsesFallback(t *testing.T) {
	result := timeout.Default(0, 5*time.Second)
	if result != 5*time.Second {
		t.Fatalf("expected 5s, got %v", result)
	}
}

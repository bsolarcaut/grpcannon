package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/grpcannon/internal/ratelimiter"
)

func TestNew_ZeroRPSReturnsNil(t *testing.T) {
	rl := ratelimiter.New(0)
	if rl != nil {
		t.Fatal("expected nil RateLimiter for zero rps")
	}
}

func TestNew_NegativeRPSReturnsNil(t *testing.T) {
	rl := ratelimiter.New(-5)
	if rl != nil {
		t.Fatal("expected nil RateLimiter for negative rps")
	}
}

func TestWait_NilLimiterReturnsImmediately(t *testing.T) {
	var rl *ratelimiter.RateLimiter
	ctx := context.Background()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWait_PermitsRequestWithinReasonableTime(t *testing.T) {
	rl := ratelimiter.New(1000) // 1000 rps → 1ms interval
	defer rl.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("expected token within timeout, got %v", err)
	}
}

func TestWait_ReturnsContextErrorOnCancellation(t *testing.T) {
	rl := ratelimiter.New(1) // 1 rps → very slow
	defer rl.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := rl.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestStop_NilLimiterDoesNotPanic(t *testing.T) {
	var rl *ratelimiter.RateLimiter
	rl.Stop() // must not panic
}

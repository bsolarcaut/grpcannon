package warmup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ashtishad/grpcannon/internal/warmup"
)

type mockCaller struct {
	callCount atomic.Int64
	errAfter  int // return error for calls > errAfter (0 = never)
	err       error
}

func (m *mockCaller) Call(_ context.Context) error {
	n := int(m.callCount.Add(1))
	if m.errAfter > 0 && n > m.errAfter {
		return m.err
	}
	return nil
}

func TestRun_ZeroRequestsIsNoop(t *testing.T) {
	caller := &mockCaller{}
	err := warmup.Run(context.Background(), caller, warmup.Config{Requests: 0})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if caller.callCount.Load() != 0 {
		t.Fatal("expected no calls")
	}
}

func TestRun_AllSucceed(t *testing.T) {
	caller := &mockCaller{}
	cfg := warmup.Config{Requests: 5, Timeout: time.Second}
	if err := warmup.Run(context.Background(), caller, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if caller.callCount.Load() != 5 {
		t.Fatalf("expected 5 calls, got %d", caller.callCount.Load())
	}
}

func TestRun_MajorityFailsReturnsError(t *testing.T) {
	caller := &mockCaller{errAfter: 1, err: errors.New("rpc error")}
	cfg := warmup.Config{Requests: 6, Timeout: time.Second}
	err := warmup.Run(context.Background(), caller, cfg)
	if err == nil {
		t.Fatal("expected error when majority fail")
	}
}

func TestRun_MinorityFailsIsOK(t *testing.T) {
	// only 1 out of 6 fails — below the half threshold
	caller := &mockCaller{errAfter: 5, err: errors.New("rpc error")}
	cfg := warmup.Config{Requests: 6, Timeout: time.Second}
	if err := warmup.Run(context.Background(), caller, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	caller := &mockCaller{}
	cfg := warmup.Config{Requests: 5, Timeout: time.Second}
	// Should not panic; errors are counted but context is already done.
	_ = warmup.Run(ctx, caller, cfg)
}

func TestDefaultConfig(t *testing.T) {
	cfg := warmup.DefaultConfig()
	if cfg.Requests != 10 {
		t.Errorf("expected 10 requests, got %d", cfg.Requests)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
}

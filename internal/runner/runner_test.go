package runner_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/grpcannon/internal/config"
	"github.com/yourusername/grpcannon/internal/runner"
)

type mockCaller struct {
	callCount atomic.Int64
	errToReturn error
	delay       time.Duration
}

func (m *mockCaller) Call(ctx context.Context) error {
	m.callCount.Add(1)
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.errToReturn
}

func baseConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Target = "localhost:50051"
	cfg.Method = "pkg.Service/Method"
	return cfg
}

func TestRunner_Run_TotalRequests(t *testing.T) {
	cfg := baseConfig()
	cfg.TotalRequests = 10
	cfg.Concurrency = 2

	caller := &mockCaller{}
	r := runner.New(cfg, caller)
	results := r.Run(context.Background())

	if len(results) != 10 {
		t.Errorf("expected 10 results, got %d", len(results))
	}
	if caller.callCount.Load() != 10 {
		t.Errorf("expected 10 calls, got %d", caller.callCount.Load())
	}
}

func TestRunner_Run_ErrorPropagated(t *testing.T) {
	cfg := baseConfig()
	cfg.TotalRequests = 5
	cfg.Concurrency = 1

	expected := errors.New("rpc error")
	caller := &mockCaller{errToReturn: expected}
	r := runner.New(cfg, caller)
	results := r.Run(context.Background())

	for _, res := range results {
		if !errors.Is(res.Err, expected) {
			t.Errorf("expected rpc error, got %v", res.Err)
		}
	}
}

func TestRunner_Run_ContextCancellation(t *testing.T) {
	cfg := baseConfig()
	cfg.TotalRequests = 100
	cfg.Concurrency = 2

	caller := &mockCaller{delay: 10 * time.Millisecond}
	r := runner.New(cfg, caller)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	results := r.Run(ctx)
	if len(results) >= 100 {
		t.Error("expected fewer than 100 results due to cancellation")
	}
}

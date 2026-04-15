package worker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/yourorg/grpcannon/internal/worker"
)

func noopCall(_ context.Context, _ *grpc.ClientConn) error {
	return nil
}

func errCall(_ context.Context, _ *grpc.ClientConn) error {
	return errors.New("rpc error")
}

func TestWorker_Run_NRequests(t *testing.T) {
	results := make(chan worker.Result, 10)
	w := worker.New(nil, noopCall, results)

	w.Run(context.Background(), 5)
	close(results)

	count := 0
	for r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error: %v", r.Err)
		}
		count++
	}
	if count != 5 {
		t.Errorf("expected 5 results, got %d", count)
	}
}

func TestWorker_Run_PropagatesErrors(t *testing.T) {
	results := make(chan worker.Result, 3)
	w := worker.New(nil, errCall, results)

	w.Run(context.Background(), 3)
	close(results)

	for r := range results {
		if r.Err == nil {
			t.Error("expected error but got nil")
		}
	}
}

func TestWorker_Run_ContextCancellation(t *testing.T) {
	results := make(chan worker.Result, 100)

	slow := func(ctx context.Context, _ *grpc.ClientConn) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	w := worker.New(nil, slow, results)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	w.Run(ctx, 0)
	close(results)

	count := 0
	for range results {
		count++
	}
	// Should have completed only a small number of requests before cancellation.
	if count > 5 {
		t.Errorf("expected few results due to cancellation, got %d", count)
	}
}

func TestWorker_Run_DurationRecorded(t *testing.T) {
	results := make(chan worker.Result, 1)
	slow := func(ctx context.Context, _ *grpc.ClientConn) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	}

	w := worker.New(nil, slow, results)
	w.Run(context.Background(), 1)
	close(results)

	r := <-results
	if r.Duration < 5*time.Millisecond {
		t.Errorf("expected duration >= 5ms, got %v", r.Duration)
	}
}

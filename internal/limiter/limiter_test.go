package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/grpcannon/internal/limiter"
)

func TestNew_DefaultsToOne(t *testing.T) {
	l := limiter.New(0)
	if l.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", l.Cap())
	}
}

func TestNew_CapReflectsN(t *testing.T) {
	l := limiter.New(5)
	if l.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", l.Cap())
	}
}

func TestAcquireRelease_InFlight(t *testing.T) {
	l := limiter.New(3)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	if l.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", l.InFlight())
	}
	l.Release()
	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight, got %d", l.InFlight())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	l := limiter.New(1)
	ctx := context.Background()
	_ = l.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected error when limiter is full")
	}
}

func TestAcquire_ContextCancellation(t *testing.T) {
	l := limiter.New(1)
	_ = l.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := l.Acquire(ctx); err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	const workers = 20
	l := limiter.New(workers)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(ctx)
			time.Sleep(5 * time.Millisecond)
			l.Release()
		}()
	}
	wg.Wait()
	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all workers done, got %d", l.InFlight())
	}
}

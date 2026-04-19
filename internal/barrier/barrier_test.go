package barrier_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/barrier"
)

func TestNew_PanicsOnZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	barrier.New(0)
}

func TestWait_ReleasesAllGoroutines(t *testing.T) {
	const n = 5
	b := barrier.New(n)
	var wg sync.WaitGroup
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errs[idx] = b.Wait(context.Background())
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %v", i, err)
		}
	}
}

func TestWait_ContextCancellation(t *testing.T) {
	b := barrier.New(3) // needs 3 but only 1 arrives
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := b.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestWait_ReusableAfterReset(t *testing.T) {
	const n = 2
	b := barrier.New(n)

	release := func() {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = b.Wait(context.Background())
			}()
		}
		wg.Wait()
	}

	release()
	release() // second generation
}

func TestReset_UnblocksWaiters(t *testing.T) {
	b := barrier.New(5)
	done := make(chan error, 1)
	go func() {
		done <- b.Wait(context.Background())
	}()
	time.Sleep(20 * time.Millisecond)
	b.Reset()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error after reset: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("goroutine not unblocked after Reset")
	}
}

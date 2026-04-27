package backpressure_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/backpressure"
)

func TestNew_PanicsOnInvalidHigh(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for high < 1")
		}
	}()
	backpressure.New(0, 0)
}

func TestNew_PanicsWhenLowGeHigh(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for low >= high")
		}
	}()
	backpressure.New(5, 5)
}

func TestAcquireRelease_Basic(t *testing.T) {
	bp := backpressure.New(10, 5)
	ctx := context.Background()

	if err := bp.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bp.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", bp.InFlight())
	}
	bp.Release()
	if bp.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight, got %d", bp.InFlight())
	}
}

func TestRecord_TripsAndLiftsPressure(t *testing.T) {
	bp := backpressure.New(10, 5)

	bp.Record(10)
	if !bp.UnderPressure() {
		t.Fatal("expected pressure at high-water mark")
	}

	bp.Record(5)
	if bp.UnderPressure() {
		t.Fatal("expected pressure lifted at low-water mark")
	}
}

func TestAcquire_BlocksUnderPressure(t *testing.T) {
	bp := backpressure.New(2, 0)
	ctx := context.Background()

	// Fill to high-water mark via Record so Acquire will block.
	bp.Record(2)

	ready := make(chan struct{})
	done := make(chan struct{})
	go func() {
		close(ready)
		_ = bp.Acquire(ctx)
		close(done)
	}()

	<-ready
	select {
	case <-done:
		t.Fatal("Acquire should be blocked")
	case <-time.After(40 * time.Millisecond):
	}

	// Lift pressure.
	bp.Record(0)

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Acquire did not unblock after pressure lifted")
	}
}

func TestAcquire_ContextCancellation(t *testing.T) {
	bp := backpressure.New(2, 0)
	bp.Record(2)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- bp.Acquire(ctx)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected context error")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Acquire did not return after context cancel")
	}
}

func TestAcquireRelease_ConcurrentSafety(t *testing.T) {
	bp := backpressure.New(50, 10)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = bp.Acquire(ctx)
			time.Sleep(time.Millisecond)
			bp.Release()
		}()
	}
	wg.Wait()

	if bp.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all goroutines done, got %d", bp.InFlight())
	}
}

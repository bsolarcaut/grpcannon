package drain_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/drain"
)

func TestNew_DefaultTimeout(t *testing.T) {
	d := drain.New(0)
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
}

func TestAcquire_ReturnsFalseAfterDrain(t *testing.T) {
	d := drain.New(time.Second)
	go d.Drain(context.Background()) //nolint:errcheck
	time.Sleep(10 * time.Millisecond)
	if d.Acquire() {
		t.Fatal("Acquire should return false after Drain called")
	}
}

func TestDrain_WaitsForInFlight(t *testing.T) {
	d := drain.New(time.Second)

	var started sync.WaitGroup
	started.Add(1)

	go func() {
		if !d.Acquire() {
			return
		}
		started.Done()
		time.Sleep(50 * time.Millisecond)
		d.Release()
	}()

	started.Wait()

	start := time.Now()
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Fatalf("drain returned too quickly: %v", elapsed)
	}
}

func TestDrain_ReturnsOnContextCancel(t *testing.T) {
	d := drain.New(5 * time.Second)

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	// never Release — simulate stuck worker

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := d.Drain(ctx)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestDrain_NoInflight_ReturnsImmediately(t *testing.T) {
	d := drain.New(time.Second)
	start := time.Now()
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("drain took too long with no in-flight work: %v", elapsed)
	}
}

func TestAcquireRelease_Concurrent(t *testing.T) {
	d := drain.New(time.Second)
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if d.Acquire() {
				time.Sleep(time.Millisecond)
				d.Release()
			}
		}()
	}
	wg.Wait()
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

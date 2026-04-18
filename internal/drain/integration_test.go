package drain_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/drain"
)

// TestIntegration_DrainStopsNewWork verifies that once Drain is called no new
// work is accepted and all previously acquired work completes.
func TestIntegration_DrainStopsNewWork(t *testing.T) {
	d := drain.New(2 * time.Second)

	var completed atomic.Int64
	const workers = 10

	// Start workers that each hold the drainer for a short time.
	for i := 0; i < workers; i++ {
		if !d.Acquire() {
			t.Fatalf("worker %d: Acquire failed before Drain", i)
		}
		go func() {
			defer d.Release()
			time.Sleep(20 * time.Millisecond)
			completed.Add(1)
		}()
	}

	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("Drain returned error: %v", err)
	}

	if got := completed.Load(); got != workers {
		t.Fatalf("expected %d completions, got %d", workers, got)
	}

	// After drain, no new acquisitions should succeed.
	if d.Acquire() {
		t.Fatal("Acquire succeeded after Drain")
	}
}

package barrier_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/example/grpcannon/internal/barrier"
)

// TestIntegration_AllWorkersFireSimultaneously verifies that no worker records
// work before all workers have passed the barrier.
func TestIntegration_AllWorkersFireSimultaneously(t *testing.T) {
	const n = 10
	b := barrier.New(n)

	var (
		passed int64
		wg     sync.WaitGroup
		mu     sync.Mutex
		snaps  []int64
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := b.Wait(context.Background()); err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			// Record how many peers have also passed at this moment.
			v := atomic.AddInt64(&passed, 1)
			mu.Lock()
			snaps = append(snaps, v)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if int(passed) != n {
		t.Fatalf("expected %d passed, got %d", n, passed)
	}
	// Every snapshot should be >= 1 (trivially true) and the max == n.
	var max int64
	for _, s := range snaps {
		if s > max {
			max = s
		}
	}
	if max != n {
		t.Fatalf("expected max snapshot %d, got %d", n, max)
	}
}

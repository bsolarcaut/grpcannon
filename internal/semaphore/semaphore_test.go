package semaphore_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/semaphore"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := semaphore.New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
}

func TestNew_ValidCapacity(t *testing.T) {
	s, err := semaphore.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Cap() != 3 {
		t.Fatalf("expected cap 3, got %d", s.Cap())
	}
}

func TestAcquireRelease_Basic(t *testing.T) {
	s, _ := semaphore.New(2)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	if s.InUse() != 1 {
		t.Fatalf("expected 1 in use, got %d", s.InUse())
	}
	s.Release()
	if s.InUse() != 0 {
		t.Fatalf("expected 0 in use, got %d", s.InUse())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	s, _ := semaphore.New(1)
	ctx := context.Background()
	_ = s.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := s.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestAcquire_ContextCancellation(t *testing.T) {
	s, _ := semaphore.New(1)
	_ = s.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := s.Acquire(ctx); err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	s, _ := semaphore.New(5)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(ctx)
			time.Sleep(5 * time.Millisecond)
			s.Release()
		}()
	}
	wg.Wait()
	if s.InUse() != 0 {
		t.Fatalf("expected 0 in use after all goroutines done, got %d", s.InUse())
	}
}

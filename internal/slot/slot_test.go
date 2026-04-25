package slot

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_CapClampedToOne(t *testing.T) {
	p := New(0)
	if p.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", p.Cap())
	}
}

func TestNew_CapReflectsN(t *testing.T) {
	p := New(5)
	if p.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", p.Cap())
	}
	if p.Available() != 5 {
		t.Fatalf("expected 5 available, got %d", p.Available())
	}
}

func TestAcquireRelease_Basic(t *testing.T) {
	p := New(3)
	ctx := context.Background()

	s, err := p.Acquire(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s < 1 || int(s) > p.Cap() {
		t.Fatalf("slot %d out of range [1, %d]", s, p.Cap())
	}
	if p.Available() != 2 {
		t.Fatalf("expected 2 available after acquire, got %d", p.Available())
	}
	if err := p.Release(s); err != nil {
		t.Fatalf("unexpected release error: %v", err)
	}
	if p.Available() != 3 {
		t.Fatalf("expected 3 available after release, got %d", p.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	p := New(1)
	ctx := context.Background()

	s, _ := p.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := p.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	_ = p.Release(s)
}

func TestAcquire_ContextCancellation(t *testing.T) {
	p := New(1)
	ctx := context.Background()
	s, _ := p.Acquire(ctx)
	defer p.Release(s) //nolint:errcheck

	ctx2, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := p.Acquire(ctx2)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestRelease_InvalidSlot(t *testing.T) {
	p := New(2)
	if err := p.Release(0); err == nil {
		t.Fatal("expected error releasing slot 0")
	}
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	const workers = 8
	p := New(workers)
	var wg sync.WaitGroup
	for i := 0; i < workers*3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := p.Acquire(context.Background())
			if err != nil {
				return
			}
			time.Sleep(time.Microsecond)
			_ = p.Release(s)
		}()
	}
	wg.Wait()
	if p.Available() != workers {
		t.Fatalf("expected all %d slots returned, got %d available", workers, p.Available())
	}
}

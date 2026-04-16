package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/throttle"
)

func TestNew_DefaultsToOne(t *testing.T) {
	th := throttle.New(0)
	if th.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", th.Cap())
	}
}

func TestNew_CapReflectsN(t *testing.T) {
	th := throttle.New(5)
	if th.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", th.Cap())
	}
}

func TestAcquireRelease_Basic(t *testing.T) {
	th := throttle.New(2)
	ctx := context.Background()

	if err := th.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	if th.Len() != 1 {
		t.Fatalf("expected len 1, got %d", th.Len())
	}
	th.Release()
	if th.Len() != 0 {
		t.Fatalf("expected len 0, got %d", th.Len())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	th := throttle.New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = th.Acquire(context.Background()) // fill the slot

	err := th.Acquire(ctx)
	if err != throttle.ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	const workers = 20
	const cap = 5
	th := throttle.New(cap)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := th.Acquire(ctx); err != nil {
				return
			}
			time.Sleep(5 * time.Millisecond)
			th.Release()
		}()
	}
	wg.Wait()
	if th.Len() != 0 {
		t.Fatalf("expected all slots released, len=%d", th.Len())
	}
}

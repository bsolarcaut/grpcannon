package shedder_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/grpcannon/internal/shedder"
)

// TestIntegration_ShedderUnderConcurrentLoad verifies that the shedder
// correctly sheds requests under concurrent load when errors pile up, and
// recovers after the cooldown elapses.
func TestIntegration_ShedderUnderConcurrentLoad(t *testing.T) {
	cfg := shedder.Config{
		Threshold:  0.5,
		WindowSize: 20,
		CoolDown:   20 * time.Millisecond,
	}
	s := shedder.New(cfg)

	// Phase 1: flood with errors to trip the shedder.
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Record(true)
		}()
	}
	wg.Wait()

	var shed int64
	for i := 0; i < 50; i++ {
		if err := s.Allow(); err == shedder.ErrShed {
			atomic.AddInt64(&shed, 1)
		}
	}
	if atomic.LoadInt64(&shed) == 0 {
		t.Fatal("expected some requests to be shed after error flood")
	}

	// Phase 2: wait for cooldown and verify recovery.
	time.Sleep(40 * time.Millisecond)
	if err := s.Allow(); err != nil {
		t.Fatalf("expected recovery after cooldown, got %v", err)
	}
}

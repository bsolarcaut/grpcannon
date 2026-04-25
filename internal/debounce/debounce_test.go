package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/debounce"
)

func TestNew_ClampsBelowZero(t *testing.T) {
	called := make(chan struct{}, 1)
	d := debounce.New(-1, func() { called <- struct{}{} })
	d.Trigger()
	select {
	case <-called:
		// ok — fired within clamped 1ms window
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected fn to be called")
	}
}

func TestTrigger_CallsFnAfterQuietPeriod(t *testing.T) {
	var count int32
	d := debounce.New(30*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	d.Trigger()
	time.Sleep(60 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}
}

func TestTrigger_ResetsTimer(t *testing.T) {
	var count int32
	d := debounce.New(40*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	// Rapid triggers should coalesce into one call.
	for i := 0; i < 5; i++ {
		d.Trigger()
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 coalesced call, got %d", got)
	}
}

func TestFlush_FiresImmediately(t *testing.T) {
	var count int32
	d := debounce.New(500*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	d.Trigger()
	d.Flush()
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call after Flush, got %d", got)
	}
}

func TestFlush_NoopWhenNoPending(t *testing.T) {
	var count int32
	d := debounce.New(10*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	// No Trigger — Flush should be a no-op.
	d.Flush()
	if got := atomic.LoadInt32(&count); got != 0 {
		t.Fatalf("expected 0 calls, got %d", got)
	}
}

func TestStop_CancelsPending(t *testing.T) {
	var count int32
	d := debounce.New(30*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	d.Trigger()
	d.Stop()
	time.Sleep(60 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 0 {
		t.Fatalf("expected 0 calls after Stop, got %d", got)
	}
}

func TestStop_Idempotent(t *testing.T) {
	d := debounce.New(10*time.Millisecond, func() {})
	d.Stop()
	d.Stop() // must not panic
}

package ticker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_DefaultInterval(t *testing.T) {
	tk := New(0, func() {})
	if tk.Interval() != time.Second {
		t.Fatalf("expected 1s default, got %v", tk.Interval())
	}
}

func TestNew_CustomInterval(t *testing.T) {
	tk := New(200*time.Millisecond, func() {})
	if tk.Interval() != 200*time.Millisecond {
		t.Fatalf("unexpected interval %v", tk.Interval())
	}
}

func TestTicker_FiresCallback(t *testing.T) {
	var count int64
	tk := New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	tk.Start(ctx)
	<-ctx.Done()
	got := atomic.LoadInt64(&count)
	if got < 3 {
		t.Fatalf("expected at least 3 ticks, got %d", got)
	}
}

func TestTicker_StopsOnStop(t *testing.T) {
	var count int64
	tk := New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	ctx := context.Background()
	tk.Start(ctx)
	time.Sleep(60 * time.Millisecond)
	tk.Stop()
	snap := atomic.LoadInt64(&count)
	time.Sleep(60 * time.Millisecond)
	if after := atomic.LoadInt64(&count); after != snap {
		t.Fatalf("ticker kept firing after Stop: %d -> %d", snap, after)
	}
}

func TestTicker_StopIdempotent(t *testing.T) {
	tk := New(50*time.Millisecond, func() {})
	tk.Start(context.Background())
	tk.Stop()
	tk.Stop() // must not panic
}

func TestTicker_StopsOnContextCancel(t *testing.T) {
	var count int64
	tk := New(20*time.Millisecond, func() {
		atomic.AddInt64(&count, 1)
	})
	ctx, cancel := context.WithCancel(context.Background())
	tk.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	cancel()
	snap := atomic.LoadInt64(&count)
	time.Sleep(60 * time.Millisecond)
	if after := atomic.LoadInt64(&count); after != snap {
		t.Fatalf("ticker kept firing after cancel: %d -> %d", snap, after)
	}
}

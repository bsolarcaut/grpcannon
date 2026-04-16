package window

import (
	"testing"
	"time"
)

func TestNew_DefaultsSize(t *testing.T) {
	w := New(0, time.Second)
	if w.size != 1 {
		t.Fatalf("expected size 1, got %d", w.size)
	}
}

func TestAdd_IncrementsCurrent(t *testing.T) {
	w := New(5, time.Second)
	w.Add(3)
	w.Add(2)
	if got := w.Sum(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestSum_EmptyWindow(t *testing.T) {
	w := New(5, time.Second)
	if got := w.Sum(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRate_ZeroWindow(t *testing.T) {
	w := New(1, 0)
	w.Add(10)
	if r := w.Rate(); r != 0 {
		t.Fatalf("expected 0 rate for zero interval, got %f", r)
	}
}

func TestRate_BasicCalculation(t *testing.T) {
	w := New(10, time.Second)
	w.Add(100)
	r := w.Rate()
	// 100 reqs / 10 sec = 10 rps
	if r < 9.9 || r > 10.1 {
		t.Fatalf("expected ~10 rps, got %f", r)
	}
}

func TestAdvance_ClearsOldBuckets(t *testing.T) {
	w := New(3, 50*time.Millisecond)
	w.Add(99)
	// Simulate time passing beyond the full window
	w.mu.Lock()
	w.last = w.last.Add(-200 * time.Millisecond)
	w.mu.Unlock()
	if got := w.Sum(); got != 0 {
		t.Fatalf("expected 0 after full rotation, got %d", got)
	}
}

func TestAdvance_PartialRotation(t *testing.T) {
	w := New(4, 50*time.Millisecond)
	w.Add(10)
	// Rotate by 2 buckets
	w.mu.Lock()
	w.last = w.last.Add(-100 * time.Millisecond)
	w.mu.Unlock()
	w.Add(5)
	if got := w.Sum(); got != 15 {
		t.Fatalf("expected 15, got %d", got)
	}
}

func TestConcurrentAdd(t *testing.T) {
	w := New(5, time.Second)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			w.Add(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if got := w.Sum(); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}

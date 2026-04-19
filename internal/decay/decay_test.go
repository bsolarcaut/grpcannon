package decay

import (
	"testing"
	"time"
)

func TestNew_DefaultHalfLifeOnZero(t *testing.T) {
	e := New(0)
	if e.alpha <= 0 || e.alpha >= 1 {
		t.Fatalf("expected alpha in (0,1), got %v", e.alpha)
	}
}

func TestAdd_FirstSampleSeeds(t *testing.T) {
	e := New(time.Second)
	e.Add(10)
	if got := e.Value(); got != 10 {
		t.Fatalf("expected 10, got %v", got)
	}
}

func TestAdd_SubsequentSampleSmooths(t *testing.T) {
	e := New(time.Second)
	e.Add(100)
	e.Add(0)
	v := e.Value()
	if v >= 100 || v <= 0 {
		t.Fatalf("expected smoothed value between 0 and 100, got %v", v)
	}
}

func TestAdd_ConvergesOnConstant(t *testing.T) {
	e := New(10 * time.Millisecond)
	for i := 0; i < 200; i++ {
		e.Add(50)
	}
	v := e.Value()
	if v < 49.9 || v > 50.1 {
		t.Fatalf("expected ~50 after convergence, got %v", v)
	}
}

func TestReset_ClearsState(t *testing.T) {
	e := New(time.Second)
	e.Add(99)
	e.Reset()
	if e.seeded {
		t.Fatal("expected seeded=false after Reset")
	}
	if e.Value() != 0 {
		t.Fatalf("expected 0 after Reset, got %v", e.Value())
	}
}

func TestAdd_AfReSeedsCorrectly(t *testing.T) {
	e := New(time.Second)
	e.Add(10)
	e.Reset()
	e.Add(77)
	if got := e.Value(); got != 77 {
		t.Fatalf("expected 77 after re-seed, got %v", got)
	}
}

func TestValue_ConcurrentSafe(t *testing.T) {
	e := New(100 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			e.Add(float64(i))
		}
		close(done)
	}()
	for i := 0; i < 500; i++ {
		_ = e.Value()
	}
	<-done
}

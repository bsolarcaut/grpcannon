package pressure_test

import (
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/pressure"
)

func TestNew_ClampsCapacity(t *testing.T) {
	g := pressure.New(0)
	if g.Capacity() != 1 {
		t.Fatalf("expected capacity 1, got %d", g.Capacity())
	}
}

func TestNew_SetsCapacity(t *testing.T) {
	g := pressure.New(10)
	if g.Capacity() != 10 {
		t.Fatalf("expected capacity 10, got %d", g.Capacity())
	}
}

func TestScore_InitiallyZero(t *testing.T) {
	g := pressure.New(5)
	if g.Score() != 0.0 {
		t.Fatalf("expected 0.0, got %f", g.Score())
	}
}

func TestAcquire_IncrementsInFlight(t *testing.T) {
	g := pressure.New(5)
	if !g.Acquire() {
		t.Fatal("expected acquire to succeed")
	}
	if g.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", g.InFlight())
	}
}

func TestAcquire_ReturnsFalseAtCapacity(t *testing.T) {
	g := pressure.New(2)
	g.Acquire()
	g.Acquire()
	if g.Acquire() {
		t.Fatal("expected acquire to fail at capacity")
	}
}

func TestRelease_DecrementsInFlight(t *testing.T) {
	g := pressure.New(5)
	g.Acquire()
	g.Release()
	if g.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight, got %d", g.InFlight())
	}
}

func TestRelease_NoOpBelowZero(t *testing.T) {
	g := pressure.New(5)
	g.Release() // should not panic or go negative
	if g.InFlight() != 0 {
		t.Fatalf("expected 0, got %d", g.InFlight())
	}
}

func TestScore_AtCapacity(t *testing.T) {
	g := pressure.New(2)
	g.Acquire()
	g.Acquire()
	if g.Score() != 1.0 {
		t.Fatalf("expected 1.0, got %f", g.Score())
	}
}

func TestScore_Partial(t *testing.T) {
	g := pressure.New(4)
	g.Acquire()
	g.Acquire()
	if got := g.Score(); got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

func TestAcquireRelease_ConcurrentSafety(t *testing.T) {
	g := pressure.New(100)
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if g.Acquire() {
				g.Release()
			}
		}()
	}
	wg.Wait()
	if g.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all releases, got %d", g.InFlight())
	}
}

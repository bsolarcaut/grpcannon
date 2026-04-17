package counter_test

import (
	"sync"
	"testing"

	"github.com/yourusername/grpcannon/internal/counter"
)

func TestNew_ZeroValues(t *testing.T) {
	c := counter.New()
	if c.Total() != 0 || c.Success() != 0 || c.Failures() != 0 || c.InFlight() != 0 {
		t.Fatal("expected all counters to start at zero")
	}
}

func TestIncTotal(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncTotal()
	if c.Total() != 2 {
		t.Fatalf("expected 2, got %d", c.Total())
	}
}

func TestIncSuccess(t *testing.T) {
	c := counter.New()
	c.IncSuccess()
	if c.Success() != 1 {
		t.Fatalf("expected 1, got %d", c.Success())
	}
}

func TestIncFailure(t *testing.T) {
	c := counter.New()
	c.IncFailure()
	c.IncFailure()
	c.IncFailure()
	if c.Failures() != 3 {
		t.Fatalf("expected 3, got %d", c.Failures())
	}
}

func TestInFlight_IncDec(t *testing.T) {
	c := counter.New()
	c.IncInFlight()
	c.IncInFlight()
	if c.InFlight() != 2 {
		t.Fatalf("expected 2 in-flight, got %d", c.InFlight())
	}
	c.DecInFlight()
	if c.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight after dec, got %d", c.InFlight())
	}
}

func TestReset(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncSuccess()
	c.IncFailure()
	c.IncInFlight()
	c.Reset()
	if c.Total() != 0 || c.Success() != 0 || c.Failures() != 0 || c.InFlight() != 0 {
		t.Fatal("expected all counters to be zero after reset")
	}
}

func TestCounter_ConcurrentSafety(t *testing.T) {
	c := counter.New()
	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.IncTotal()
			c.IncSuccess()
			c.IncInFlight()
			c.DecInFlight()
		}()
	}
	wg.Wait()
	if c.Total() != goroutines {
		t.Fatalf("expected %d total, got %d", goroutines, c.Total())
	}
	if c.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all done, got %d", c.InFlight())
	}
}

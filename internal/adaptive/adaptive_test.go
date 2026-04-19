package adaptive

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	c := New(Config{})
	if c.Current() != defaultMin {
		t.Fatalf("expected %d, got %d", defaultMin, c.Current())
	}
	if c.max != defaultMax {
		t.Fatalf("expected max %d, got %d", defaultMax, c.max)
	}
}

func TestNew_CustomInitial(t *testing.T) {
	c := New(Config{Initial: 10, Min: 1, Max: 100})
	if c.Current() != 10 {
		t.Fatalf("expected 10, got %d", c.Current())
	}
}

func TestScaleUp_IncreasesBy_Step(t *testing.T) {
	c := New(Config{Initial: 4, Min: 1, Max: 100, Step: 3, Cooldown: 0})
	n := c.ScaleUp()
	if n != 7 {
		t.Fatalf("expected 7, got %d", n)
	}
}

func TestScaleDown_DecreasesBy_Step(t *testing.T) {
	c := New(Config{Initial: 10, Min: 1, Max: 100, Step: 3, Cooldown: 0})
	n := c.ScaleDown()
	if n != 7 {
		t.Fatalf("expected 7, got %d", n)
	}
}

func TestScaleUp_CapsAtMax(t *testing.T) {
	c := New(Config{Initial: 9, Min: 1, Max: 10, Step: 5, Cooldown: 0})
	n := c.ScaleUp()
	if n != 10 {
		t.Fatalf("expected 10, got %d", n)
	}
}

func TestScaleDown_FloorsAtMin(t *testing.T) {
	c := New(Config{Initial: 2, Min: 1, Max: 10, Step: 5, Cooldown: 0})
	n := c.ScaleDown()
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
}

func TestScaleUp_RespectsCoooldown(t *testing.T) {
	c := New(Config{Initial: 4, Min: 1, Max: 100, Step: 2, Cooldown: 10 * time.Second})
	c.ScaleUp() // first call sets lastAdj
	n := c.ScaleUp() // should be blocked by cooldown
	if n != 6 {
		t.Fatalf("expected 6 (no second scale), got %d", n)
	}
}

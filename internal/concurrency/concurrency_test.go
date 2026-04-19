package concurrency_test

import (
	"testing"

	"github.com/lukeberry99/grpcannon/internal/concurrency"
)

func TestNew_ClampsMin(t *testing.T) {
	c := concurrency.New(0, 0, 10, 1)
	if c.Current() < 1 {
		t.Fatalf("expected current >= 1, got %d", c.Current())
	}
}

func TestNew_InitialWithinBounds(t *testing.T) {
	c := concurrency.New(5, 1, 10, 1)
	if got := c.Current(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestNew_ClampsInitialToMax(t *testing.T) {
	c := concurrency.New(20, 1, 10, 1)
	if got := c.Current(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestNew_ClampsInitialToMin(t *testing.T) {
	c := concurrency.New(0, 3, 10, 1)
	if got := c.Current(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestScaleUp_IncreasesLevel(t *testing.T) {
	c := concurrency.New(5, 1, 10, 2)
	got := c.ScaleUp()
	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestScaleUp_CapsAtMax(t *testing.T) {
	c := concurrency.New(10, 1, 10, 2)
	got := c.ScaleUp()
	if got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestScaleDown_DecreasesLevel(t *testing.T) {
	c := concurrency.New(5, 1, 10, 2)
	got := c.ScaleDown()
	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestScaleDown_FloorsAtMin(t *testing.T) {
	c := concurrency.New(1, 1, 10, 5)
	got := c.ScaleDown()
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestHistory_RecordsChanges(t *testing.T) {
	c := concurrency.New(5, 1, 20, 1)
	c.ScaleUp()
	c.ScaleDown()
	h := c.History()
	if len(h) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(h))
	}
}

func TestHistory_ReturnsCopy(t *testing.T) {
	c := concurrency.New(5, 1, 20, 1)
	c.ScaleUp()
	h1 := c.History()
	h1[0].Workers = 999
	h2 := c.History()
	if h2[0].Workers == 999 {
		t.Fatal("history slice is not a copy")
	}
}

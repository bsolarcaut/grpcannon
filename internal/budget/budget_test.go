package budget_test

import (
	"errors"
	"testing"

	"github.com/example/grpcannon/internal/budget"
)

var errFake = errors.New("rpc error")

func TestNew_Defaults(t *testing.T) {
	b := budget.New(0, 0)
	if b == nil {
		t.Fatal("expected non-nil budget")
	}
}

func TestCheck_BelowMinTotal_NeverTrips(t *testing.T) {
	b := budget.New(0.01, 20)
	for i := 0; i < 19; i++ {
		b.Record(errFake)
	}
	if err := b.Check(); err != nil {
		t.Fatalf("expected nil before minTotal, got %v", err)
	}
}

func TestCheck_BelowThreshold_OK(t *testing.T) {
	b := budget.New(0.10, 10)
	for i := 0; i < 9; i++ {
		b.Record(nil)
	}
	b.Record(errFake) // 10 % failures
	if err := b.Check(); err != nil {
		t.Fatalf("expected nil at threshold boundary, got %v", err)
	}
}

func TestCheck_ExceedsThreshold_Exhausted(t *testing.T) {
	b := budget.New(0.10, 10)
	for i := 0; i < 8; i++ {
		b.Record(nil)
	}
	b.Record(errFake)
	b.Record(errFake) // 20 % failures
	if err := b.Check(); !errors.Is(err, budget.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
}

func TestRate_Empty(t *testing.T) {
	b := budget.New(0.05, 10)
	if r := b.Rate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRate_Partial(t *testing.T) {
	b := budget.New(0.05, 1)
	b.Record(nil)
	b.Record(errFake)
	if r := b.Rate(); r != 0.5 {
		t.Fatalf("expected 0.5, got %f", r)
	}
}

func TestSnapshot(t *testing.T) {
	b := budget.New(0.05, 1)
	b.Record(nil)
	b.Record(errFake)
	total, failures := b.Snapshot()
	if total != 2 || failures != 1 {
		t.Fatalf("expected 2/1, got %d/%d", total, failures)
	}
}

func TestCheck_AllSuccess(t *testing.T) {
	b := budget.New(0.05, 5)
	for i := 0; i < 10; i++ {
		b.Record(nil)
	}
	if err := b.Check(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

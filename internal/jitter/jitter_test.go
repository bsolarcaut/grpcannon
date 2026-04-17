package jitter_test

import (
	"testing"
	"time"

	"github.com/gkampitakis/grpcannon/internal/jitter"
)

func TestFull_ZeroReturnsZero(t *testing.T) {
	j := jitter.New(42)
	if got := j.Full(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFull_NegativeReturnsZero(t *testing.T) {
	j := jitter.New(42)
	if got := j.Full(-time.Second); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFull_WithinRange(t *testing.T) {
	j := jitter.New(1)
	base := 100 * time.Millisecond
	for i := 0; i < 200; i++ {
		v := j.Full(base)
		if v < 0 || v >= base {
			t.Fatalf("Full(%v) = %v out of [0, %v)", base, v, base)
		}
	}
}

func TestEqual_ZeroReturnsZero(t *testing.T) {
	j := jitter.New(42)
	if got := j.Equal(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestEqual_WithinHalfRange(t *testing.T) {
	j := jitter.New(7)
	base := 200 * time.Millisecond
	for i := 0; i < 200; i++ {
		v := j.Equal(base)
		if v < base/2 || v >= base {
			t.Fatalf("Equal(%v) = %v out of [%v, %v)", base, v, base/2, base)
		}
	}
}

func TestApply_ClampsToMin(t *testing.T) {
	j := jitter.New(99)
	// base=2ms so Equal returns [1ms,2ms); min=5ms should clamp up
	v := j.Apply(2*time.Millisecond, 5*time.Millisecond, time.Second)
	if v != 5*time.Millisecond {
		t.Fatalf("expected min clamp 5ms, got %v", v)
	}
}

func TestApply_ClampsToMax(t *testing.T) {
	j := jitter.New(99)
	v := j.Apply(time.Second, 0, 10*time.Millisecond)
	if v != 10*time.Millisecond {
		t.Fatalf("expected max clamp 10ms, got %v", v)
	}
}

func TestApply_WithinBounds(t *testing.T) {
	j := jitter.New(3)
	min, max := 50*time.Millisecond, 500*time.Millisecond
	for i := 0; i < 100; i++ {
		v := j.Apply(200*time.Millisecond, min, max)
		if v < min || v > max {
			t.Fatalf("Apply result %v out of [%v, %v]", v, min, max)
		}
	}
}

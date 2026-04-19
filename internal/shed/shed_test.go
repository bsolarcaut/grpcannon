package shed

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	s := New(0.5, 0, 0)
	if s.size != 100 {
		t.Fatalf("expected default window 100, got %d", s.size)
	}
	if s.cooldown != 100*time.Millisecond {
		t.Fatalf("expected default cooldown 100ms, got %v", s.cooldown)
	}
}

func TestAllow_EmptyWindowPermits(t *testing.T) {
	s := New(0.5, 10, time.Millisecond)
	if !s.Allow() {
		t.Fatal("empty window should always allow")
	}
}

func TestRate_AllSuccess(t *testing.T) {
	s := New(0.5, 10, time.Millisecond)
	for i := 0; i < 10; i++ {
		s.Record(false)
	}
	if s.Rate() != 0 {
		t.Fatalf("expected 0 error rate, got %f", s.Rate())
	}
}

func TestRate_AllErrors(t *testing.T) {
	s := New(0.5, 10, time.Millisecond)
	for i := 0; i < 10; i++ {
		s.Record(true)
	}
	if s.Rate() != 1.0 {
		t.Fatalf("expected 1.0 error rate, got %f", s.Rate())
	}
}

func TestAllow_BelowThreshold(t *testing.T) {
	s := New(0.5, 10, time.Millisecond)
	for i := 0; i < 10; i++ {
		s.Record(false)
	}
	s.Record(true) // 1/10 = 0.1 < 0.5
	if !s.Allow() {
		t.Fatal("rate below threshold should allow")
	}
}

func TestAllow_AboveThreshold_Sheds(t *testing.T) {
	s := New(0.5, 10, time.Millisecond)
	for i := 0; i < 10; i++ {
		s.Record(true)
	}
	time.Sleep(2 * time.Millisecond) // pass cooldown
	if s.Allow() {
		t.Fatal("rate above threshold should shed")
	}
}

func TestAllow_CooldownPreventsRepeatShed(t *testing.T) {
	s := New(0.5, 10, 500*time.Millisecond)
	for i := 0; i < 10; i++ {
		s.Record(true)
	}
	// First shed decision sets lastShed
	s.Allow()
	// Second call within cooldown should also return false
	if s.Allow() {
		t.Fatal("within cooldown should still shed")
	}
}

func TestRecord_WindowWraps(t *testing.T) {
	s := New(0.5, 4, time.Millisecond)
	// Fill with errors
	for i := 0; i < 4; i++ {
		s.Record(true)
	}
	// Overwrite with successes
	for i := 0; i < 4; i++ {
		s.Record(false)
	}
	if s.Rate() != 0 {
		t.Fatalf("window should have wrapped to all success, got %f", s.Rate())
	}
}

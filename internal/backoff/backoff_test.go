package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/backoff"
)

func TestExponential_Wait_GrowsExponentially(t *testing.T) {
	e := &backoff.Exponential{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
	}
	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
	}
	for i, want := range expected {
		got := e.Wait(i)
		if got != want {
			t.Errorf("attempt %d: got %v, want %v", i, got, want)
		}
	}
}

func TestExponential_Wait_CapsAtMax(t *testing.T) {
	e := &backoff.Exponential{
		Base:       1 * time.Second,
		Max:        2 * time.Second,
		Multiplier: 4.0,
	}
	got := e.Wait(5)
	if got != 2*time.Second {
		t.Errorf("expected max 2s, got %v", got)
	}
}

func TestExponential_Wait_DefaultMultiplier(t *testing.T) {
	e := &backoff.Exponential{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 0, // should default to 2.0
	}
	got := e.Wait(1)
	if got != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", got)
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := backoff.Do(context.Background(), backoff.DefaultExponential(), 3, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndEventuallySucceeds(t *testing.T) {
	calls := 0
	err := backoff.Do(context.Background(), backoff.DefaultExponential(), 5, func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsLastErrorAfterMaxAttempts(t *testing.T) {
	sentinel := errors.New("always fails")
	err := backoff.Do(context.Background(), backoff.DefaultExponential(), 3, func() error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := backoff.Do(ctx, backoff.DefaultExponential(), 5, func() error {
		calls++
		return errors.New("err")
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Errorf("expected 0 calls after pre-cancelled ctx, got %d", calls)
	}
}

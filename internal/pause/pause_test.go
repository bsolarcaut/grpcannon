package pause

import (
	"context"
	"testing"
	"time"
)

func TestNew_NotPaused(t *testing.T) {
	c := New()
	if c.IsPaused() {
		t.Fatal("expected controller to start unpaused")
	}
}

func TestPause_SetsPaused(t *testing.T) {
	c := New()
	c.Pause()
	if !c.IsPaused() {
		t.Fatal("expected controller to be paused")
	}
}

func TestResume_ClearsPaused(t *testing.T) {
	c := New()
	c.Pause()
	c.Resume()
	if c.IsPaused() {
		t.Fatal("expected controller to be unpaused after Resume")
	}
}

func TestWait_ReturnsImmediatelyWhenNotPaused(t *testing.T) {
	c := New()
	done := make(chan error, 1)
	go func() { done <- c.Wait(context.Background()) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait blocked unexpectedly")
	}
}

func TestWait_BlocksWhenPaused(t *testing.T) {
	c := New()
	c.Pause()
	done := make(chan error, 1)
	go func() { done <- c.Wait(context.Background()) }()
	select {
	case <-done:
		t.Fatal("Wait returned while paused")
	case <-time.After(50 * time.Millisecond):
	}
	c.Resume()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error after resume: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait did not unblock after Resume")
	}
}

func TestWait_ContextCancellation(t *testing.T) {
	c := New()
	c.Pause()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- c.Wait(ctx) }()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected context error, got nil")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait did not unblock after context cancel")
	}
}

func TestStop_UnblocksWaiters(t *testing.T) {
	c := New()
	c.Pause()
	done := make(chan error, 1)
	go func() { done <- c.Wait(context.Background()) }()
	c.Stop()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error after stop: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait did not unblock after Stop")
	}
}

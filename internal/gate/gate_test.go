package gate_test

import (
	"sync"
	"testing"
	"time"

	"github.com/romangurevitch/grpcannon/internal/gate"
)

func TestNew_IsOpen(t *testing.T) {
	g := gate.New()
	if !g.IsOpen() {
		t.Fatal("expected gate to be open after New")
	}
}

func TestPause_ClosesGate(t *testing.T) {
	g := gate.New()
	g.Pause()
	if g.IsOpen() {
		t.Fatal("expected gate to be closed after Pause")
	}
}

func TestResume_OpensGate(t *testing.T) {
	g := gate.New()
	g.Pause()
	g.Resume()
	if !g.IsOpen() {
		t.Fatal("expected gate to be open after Resume")
	}
}

func TestWait_BlocksWhenPaused(t *testing.T) {
	g := gate.New()
	g.Pause()

	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Wait returned while gate was paused")
	case <-time.After(50 * time.Millisecond):
	}

	g.Resume()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Wait did not return after Resume")
	}
}

func TestStop_UnblocksWaiters(t *testing.T) {
	g := gate.New()
	g.Pause()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Wait()
		}()
	}

	g.Stop()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Stop did not unblock all waiters")
	}
}

func TestWait_OpenGateReturnsImmediately(t *testing.T) {
	g := gate.New()
	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait blocked on open gate")
	}
}

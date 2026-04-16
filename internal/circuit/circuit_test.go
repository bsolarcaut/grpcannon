package circuit_test

import (
	"testing"
	"time"

	"github.com/dev/grpcannon/internal/circuit"
)

func TestNew_Defaults(t *testing.T) {
	b := circuit.New(0, 0)
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
	if b.CurrentState() != circuit.StateClosed {
		t.Fatalf("expected Closed, got %v", b.CurrentState())
	}
}

func TestAllow_ClosedPermits(t *testing.T) {
	b := circuit.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuit.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.CurrentState() != circuit.StateOpen {
		t.Fatalf("expected Open, got %v", b.CurrentState())
	}
	if err := b.Allow(); err != circuit.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	b := circuit.New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != circuit.StateOpen {
		t.Fatal("expected Open")
	}
	b.RecordSuccess()
	if b.CurrentState() != circuit.StateClosed {
		t.Fatalf("expected Closed after success, got %v", b.CurrentState())
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b := circuit.New(1, 10*time.Millisecond)
	b.RecordFailure()
	if err := b.Allow(); err != circuit.ErrOpen {
		t.Fatalf("expected ErrOpen immediately, got %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	if b.CurrentState() != circuit.StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", b.CurrentState())
	}
}

func TestConcurrentRecordFailure(t *testing.T) {
	b := circuit.New(100, time.Second)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() { b.RecordFailure(); done <- struct{}{} }()
		go func() { b.RecordSuccess(); done <- struct{}{} }()
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}

package circuit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dev/grpcannon/internal/circuit"
)

// call simulates a gRPC call that may fail.
func call(fail bool) error {
	if fail {
		return errors.New("rpc error")
	}
	return nil
}

func guardedCall(b *circuit.Breaker, fail bool) error {
	if err := b.Allow(); err != nil {
		return err
	}
	err := call(fail)
	if err != nil {
		b.RecordFailure()
	} else {
		b.RecordSuccess()
	}
	return err
}

func TestIntegration_CircuitOpensAndRecovers(t *testing.T) {
	b := circuit.New(3, 20*time.Millisecond)

	// Three failures open the circuit.
	for i := 0; i < 3; i++ {
		if err := guardedCall(b, true); err == nil {
			t.Fatal("expected error")
		}
	}
	if b.CurrentState() != circuit.StateOpen {
		t.Fatal("expected Open")
	}

	// Subsequent calls are rejected without invoking the RPC.
	if err := guardedCall(b, false); !errors.Is(err, circuit.ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}

	// After cooldown a successful probe closes the circuit.
	time.Sleep(30 * time.Millisecond)
	if err := guardedCall(b, false); err != nil {
		t.Fatalf("probe should succeed, got %v", err)
	}
	if b.CurrentState() != circuit.StateClosed {
		t.Fatalf("expected Closed after recovery, got %v", b.CurrentState())
	}
}

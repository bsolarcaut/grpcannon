package concurrency_test

import (
	"testing"
	"time"

	"github.com/lukeberry99/grpcannon/internal/concurrency"
)

func TestDefaultPolicy_Fields(t *testing.T) {
	p := concurrency.DefaultPolicy()
	if p.MaxErrorRate != 0.05 {
		t.Fatalf("expected MaxErrorRate 0.05, got %v", p.MaxErrorRate)
	}
	if p.MaxP99Latency != 500*time.Millisecond {
		t.Fatalf("expected MaxP99Latency 500ms, got %v", p.MaxP99Latency)
	}
}

func TestEvaluate_HighErrorRate_ScalesDown(t *testing.T) {
	p := concurrency.DefaultPolicy()
	d := p.Evaluate(0.10, 100*time.Millisecond)
	if d != concurrency.ScaleDown {
		t.Fatalf("expected ScaleDown, got %v", d)
	}
}

func TestEvaluate_HighLatency_ScalesDown(t *testing.T) {
	p := concurrency.DefaultPolicy()
	d := p.Evaluate(0.00, 600*time.Millisecond)
	if d != concurrency.ScaleDown {
		t.Fatalf("expected ScaleDown, got %v", d)
	}
}

func TestEvaluate_LowErrorRate_ScalesUp(t *testing.T) {
	p := concurrency.DefaultPolicy()
	d := p.Evaluate(0.005, 100*time.Millisecond)
	if d != concurrency.ScaleUp {
		t.Fatalf("expected ScaleUp, got %v", d)
	}
}

func TestEvaluate_MidRange_Holds(t *testing.T) {
	p := concurrency.DefaultPolicy()
	d := p.Evaluate(0.03, 200*time.Millisecond)
	if d != concurrency.Hold {
		t.Fatalf("expected Hold, got %v", d)
	}
}

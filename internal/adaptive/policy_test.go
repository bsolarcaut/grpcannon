package adaptive

import "testing"

func TestDefaultPolicy_Fields(t *testing.T) {
	p := DefaultPolicy()
	if p.ErrorThreshold <= 0 {
		t.Fatal("expected positive ErrorThreshold")
	}
	if p.LatencyThreshold <= 0 {
		t.Fatal("expected positive LatencyThreshold")
	}
}

func TestEvaluate_HighErrorRate_ScalesDown(t *testing.T) {
	p := DefaultPolicy()
	if got := p.Evaluate(0.10, 0.1); got != ScaleDown {
		t.Fatalf("expected ScaleDown, got %v", got)
	}
}

func TestEvaluate_HighLatency_ScalesDown(t *testing.T) {
	p := DefaultPolicy()
	if got := p.Evaluate(0.0, 2.0); got != ScaleDown {
		t.Fatalf("expected ScaleDown, got %v", got)
	}
}

func TestEvaluate_LowErrorRate_ScalesUp(t *testing.T) {
	p := DefaultPolicy()
	if got := p.Evaluate(0.005, 0.1); got != ScaleUp {
		t.Fatalf("expected ScaleUp, got %v", got)
	}
}

func TestEvaluate_MidRange_Holds(t *testing.T) {
	p := DefaultPolicy()
	// error rate between TargetErrorRate and ErrorThreshold
	if got := p.Evaluate(0.03, 0.1); got != Hold {
		t.Fatalf("expected Hold, got %v", got)
	}
}

func TestEvaluate_BothThresholdsExceeded_ScalesDown(t *testing.T) {
	p := DefaultPolicy()
	if got := p.Evaluate(0.9, 9.9); got != ScaleDown {
		t.Fatalf("expected ScaleDown, got %v", got)
	}
}

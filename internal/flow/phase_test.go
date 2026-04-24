package flow

import (
	"testing"
)

func TestPhase_String(t *testing.T) {
	cases := []struct {
		phase Phase
		want  string
	}{
		{PhaseIdle, "idle"},
		{PhaseWarmup, "warmup"},
		{PhaseRamp, "ramp"},
		{PhaseSteady, "steady"},
		{PhaseDrain, "drain"},
		{PhaseDone, "done"},
		{Phase(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.phase.String(); got != tc.want {
			t.Errorf("Phase(%d).String() = %q, want %q", tc.phase, got, tc.want)
		}
	}
}

func TestPhase_IsFinal(t *testing.T) {
	if PhaseIdle.IsFinal() {
		t.Error("PhaseIdle should not be final")
	}
	if PhaseRamp.IsFinal() {
		t.Error("PhaseRamp should not be final")
	}
	if !PhaseDone.IsFinal() {
		t.Error("PhaseDone should be final")
	}
}

func TestPhase_IsActive(t *testing.T) {
	cases := []struct {
		phase    Phase
		active   bool
	}{
		{PhaseIdle, false},
		{PhaseWarmup, true},
		{PhaseRamp, true},
		{PhaseSteady, true},
		{PhaseDrain, true},
		{PhaseDone, false},
	}
	for _, tc := range cases {
		if got := tc.phase.IsActive(); got != tc.active {
			t.Errorf("Phase(%s).IsActive() = %v, want %v", tc.phase, got, tc.active)
		}
	}
}

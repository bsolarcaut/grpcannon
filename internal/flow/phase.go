package flow

// Phase represents a named stage in the load-test lifecycle.
type Phase int

const (
	PhaseIdle Phase = iota
	PhaseWarmup
	PhaseRamp
	PhaseSteady
	PhaseDrain
	PhaseDone
)

var phaseNames = map[Phase]string{
	PhaseIdle:    "idle",
	PhaseWarmup:  "warmup",
	PhaseRamp:    "ramp",
	PhaseSteady:  "steady",
	PhaseDrain:   "drain",
	PhaseDone:    "done",
}

// String returns the human-readable name of the phase.
func (p Phase) String() string {
	if name, ok := phaseNames[p]; ok {
		return name
	}
	return "unknown"
}

// IsFinal reports whether the phase represents a terminal state.
func (p Phase) IsFinal() bool {
	return p == PhaseDone
}

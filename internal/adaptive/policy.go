package adaptive

// Policy decides whether to scale up, scale down, or hold based on
// observed metrics.
type Policy struct {
	// ErrorThreshold is the error rate [0,1] above which we scale down.
	ErrorThreshold float64
	// LatencyThreshold is the p99 latency above which we scale down.
	LatencyThreshold float64
	// TargetErrorRate is the error rate below which we may scale up.
	TargetErrorRate float64
}

// Direction indicates the scaling direction.
type Direction int

const (
	Hold     Direction = 0
	ScaleUp  Direction = 1
	ScaleDown Direction = -1
)

// DefaultPolicy returns a sensible default Policy.
func DefaultPolicy() Policy {
	return Policy{
		ErrorThreshold:   0.05,
		LatencyThreshold: 1.0, // seconds
		TargetErrorRate:  0.01,
	}
}

// Evaluate returns a scaling Direction given current error rate and p99 latency (seconds).
func (p Policy) Evaluate(errorRate, p99Latency float64) Direction {
	if errorRate > p.ErrorThreshold || p99Latency > p.LatencyThreshold {
		return ScaleDown
	}
	 p.TargetErrorRate {
		return ScaleUp
	}
	return Hold
}

package concurrency

import "time"

// Decision represents a scaling action.
type Decision int

const (
	Hold     Decision = iota
	ScaleUp
	ScaleDown
)

// Policy evaluates metrics and returns a scaling Decision.
type Policy struct {
	MaxErrorRate  float64
	MaxP99Latency time.Duration
	MinErrorRate  float64
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		MaxErrorRate:  0.05,
		MinErrorRate:  0.01,
		MaxP99Latency: 500 * time.Millisecond,
	}
}

// Evaluate returns a Decision given current error rate and p99 latency.
func (p Policy) Evaluate(errorRate float64, p99 time.Duration) Decision {
	if errorRate > p.MaxErrorRate || p99 > p.MaxP99Latency {
		return ScaleDown
	}
	if errorRate < p.MinErrorRate && p99 < p.MaxP99Latency {
		return ScaleUp
	}
	return Hold
}

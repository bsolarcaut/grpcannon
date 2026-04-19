// Package adaptive implements a feedback-driven concurrency controller for
// grpcannon load runs.
//
// The Controller tracks a current worker count and exposes ScaleUp / ScaleDown
// methods that are gated by a configurable cooldown period so the system
// avoids oscillating too quickly.
//
// A Policy translates observed metrics (error rate, p99 latency) into a
// scaling Direction which the caller can act on:
//
//	pol := adaptive.DefaultPolicy()
//	ctrl := adaptive.New(adaptive.Config{Min: 2, Max: 64})
//
//	switch pol.Evaluate(errRate, p99) {
//	case adaptive.ScaleUp:
//	    n := ctrl.ScaleUp()
//	    // spawn workers up to n
//	case adaptive.ScaleDown:
//	    n := ctrl.ScaleDown()
//	    // retire workers down to n
//	}
package adaptive

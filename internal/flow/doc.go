// Package flow wires together the warmup, ramp, steady-state, and drain
// phases of a grpcannon load test into a single, cancellable Run call.
//
// Typical usage:
//
//	err := flow.Run(ctx, flow.Config{
//		Warmup: warmup.DefaultConfig(callFn),
//		Ramp:   ramp.Default(),
//		Steady: 30 * time.Second,
//		Runner: r,
//	})
package flow

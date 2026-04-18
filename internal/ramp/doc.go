// Package ramp implements a step-based concurrency ramp-up schedule for
// grpcannon load tests.
//
// Instead of launching all workers at once, ramp gradually increases
// concurrency from StartWorkers to EndWorkers over a configurable Duration
// split into discrete Steps. This reduces the risk of overwhelming a target
// service at the start of a test run.
//
// Basic usage:
//
//	cfg := ramp.Default()
//	ch := ramp.Schedule(ctx, cfg)
//	for level := range ch {
//		pool.Resize(level)
//	}
package ramp

// Package throttle implements a semaphore-based concurrency limiter for
// grpcannon. It is used by the worker pool to cap the number of in-flight
// gRPC calls at any given time, preventing the target service from being
// overwhelmed beyond the configured --concurrency value.
//
// Usage:
//
//	th := throttle.New(cfg.Concurrency)
//	if err := th.Acquire(ctx); err != nil {
//		return err
//	}
//	defer th.Release()
//	// ... perform gRPC call ...
package throttle

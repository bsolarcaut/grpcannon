// Package backpressure provides a high/low-water-mark controller for
// regulating the rate at which new work is accepted under load.
//
// When the number of in-flight requests reaches the high-water mark the
// controller enters a pressured state and any caller that invokes Acquire
// will block until the count falls back to the low-water mark, at which
// point all waiters are released simultaneously.
//
// Typical usage:
//
//	bp := backpressure.New(100, 80)
//
//	if err := bp.Acquire(ctx); err != nil {
//	    return err
//	}
//	defer bp.Release()
//	// … do work …
package backpressure

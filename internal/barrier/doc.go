// Package barrier implements a cyclic synchronisation barrier.
//
// A Barrier is useful when a fixed number of concurrent workers must all reach
// a common checkpoint before any of them proceeds — for example, aligning the
// start of every load-test round so that all workers fire simultaneously.
//
// Usage:
//
//	b := barrier.New(workers)
//	// inside each worker goroutine:
//	if err := b.Wait(ctx); err != nil {
//	    return err // context cancelled
//	}
//	// all workers proceed together from here
package barrier

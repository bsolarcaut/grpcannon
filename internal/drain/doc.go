// Package drain implements a graceful-shutdown barrier for grpcannon.
//
// When a shutdown signal is received the runner stops issuing new requests
// and calls Drain to wait for any requests that are already in-flight to
// complete.  A configurable timeout prevents the process from hanging
// indefinitely if a backend is unresponsive.
//
// Typical usage:
//
//	d := drain.New(3 * time.Second)
//
//	// inside worker loop
//	if !d.Acquire() {
//		return // shutting down
//	}
//	defer d.Release()
//	// … do work …
//
//	// on shutdown signal
//	if err := d.Drain(ctx); err != nil {
//		log.Println("drain timeout:", err)
//	}
package drain

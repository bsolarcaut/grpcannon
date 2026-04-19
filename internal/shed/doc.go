// Package shed provides a simple load-shedding primitive for grpcannon.
//
// A Shed tracks the recent error rate of outgoing gRPC calls using a
// fixed-size circular window. Once the error rate rises above a configured
// threshold the Shed begins rejecting new requests via Allow, giving the
// target service time to recover.
//
// Typical usage:
//
//	s := shed.New(0.5, 200, 250*time.Millisecond)
//	if !s.Allow() {
//	    // drop / return early
//	}
//	// … perform call …
//	s.Record(err != nil)
package shed

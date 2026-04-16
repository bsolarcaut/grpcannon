// Package snapshot provides a thread-safe collector for capturing
// point-in-time metrics during a grpcannon load test run.
//
// Typical usage:
//
//	c := snapshot.NewCollector()
//	c.RecordSuccess()
//	m := c.Snap()
package snapshot

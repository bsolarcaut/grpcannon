// Package histogram implements a fixed-boundary latency histogram for
// grpcannon load-test results.
//
// Boundaries are expressed as time.Duration values. Samples are placed
// into the first bucket whose upper bound exceeds the sample; samples
// larger than all bounds fall into an overflow bucket.
//
// Usage:
//
//	h := histogram.New([]time.Duration{
//		1*time.Millisecond, 5*time.Millisecond, 10*time.Millisecond,
//	})
//	h.Record(3 * time.Millisecond)
//	h.Print(os.Stdout)
package histogram

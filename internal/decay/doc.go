// Package decay implements an exponentially weighted moving average (EWMA)
// for smoothing noisy time-series signals such as per-second error rates
// and tail latencies produced during a load test.
//
// Usage:
//
//	ew := decay.New(5 * time.Second) // half-life of 5 s
//	ew.Add(42.0)
//	fmt.Println(ew.Value())
package decay

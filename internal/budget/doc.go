// Package budget provides an error-budget guard for grpcannon load tests.
//
// An error budget defines the maximum tolerable failure rate for a test run.
// Once the observed failure ratio exceeds the configured threshold (and a
// minimum number of requests have been made), Budget.Check returns
// ErrExhausted, allowing the runner to abort early rather than continue
// hammering a degraded service.
//
// Example:
//
//	b := budget.New(0.05, 20) // 5 % threshold, at least 20 requests
//	b.Record(nil)             // success
//	b.Record(err)             // failure
//	if err := b.Check(); err != nil {
//	    // abort the run
//	}
package budget

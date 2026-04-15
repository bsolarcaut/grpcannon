package stats

import "time"

// Summary holds aggregated statistics computed from a set of Results.
type Summary struct {
	Total     int
	Succeeded int
	Failed    int
	Min       time.Duration
	Mean      time.Duration
	P50       time.Duration
	P95       time.Duration
	P99       time.Duration
	Max       time.Duration
}

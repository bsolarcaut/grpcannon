package stats

import (
	"fmt"
	"sort"
	"time"
)

// Result holds the outcome of a single RPC call.
type Result struct {
	Latency time.Duration
	Err     error
}

// Compute calculates a Summary from a slice of Results.
func Compute(results []Result) Summary {
	if len(results) == 0 {
		return Summary{}
	}

	var latencies []time.Duration
	var total, succeeded, failed int
	var sum time.Duration

	for _, r := range results {
		total++
		if r.Err != nil {
			failed++
			continue
		}
		succeeded++
		latencies = append(latencies, r.Latency)
		sum += r.Latency
	}

	if len(latencies) == 0 {
		return Summary{Total: total, Succeeded: succeeded, Failed: failed}
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	mean := sum / time.Duration(len(latencies))

	return Summary{
		Total:     total,
		Succeeded: succeeded,
		Failed:    failed,
		Min:       latencies[0],
		Mean:      mean,
		P50:       percentile(latencies, 50),
		P95:       percentile(latencies, 95),
		P99:       percentile(latencies, 99),
		Max:       latencies[len(latencies)-1],
	}
}

// percentile returns the p-th percentile duration from a sorted slice.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p / 100.0)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// Print writes a human-readable summary to stdout.
func Print(s Summary) {
	fmt.Printf("Total: %d | OK: %d | Err: %d\n", s.Total, s.Succeeded, s.Failed)
	fmt.Printf("Min: %v | Mean: %v | P50: %v | P95: %v | P99: %v | Max: %v\n",
		s.Min, s.Mean, s.P50, s.P95, s.P99, s.Max)
}

package stats

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/yourusername/grpcannon/internal/runner"
)

// Summary aggregates load test results.
type Summary struct {
	Total    int
	Success  int
	Failures int
	Mean     time.Duration
	Min      time.Duration
	Max      time.Duration
	P50      time.Duration
	P95      time.Duration
	P99      time.Duration
}

// Compute builds a Summary from a slice of runner Results.
func Compute(results []runner.Result) Summary {
	if len(results) == 0 {
		return Summary{}
	}

	durations := make([]time.Duration, 0, len(results))
	var total time.Duration
	failures := 0

	for _, r := range results {
		if r.Err != nil {
			failures++
		} else {
			durations = append(durations, r.Duration)
			total += r.Duration
		}
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	s := Summary{
		Total:    len(results),
		Failures: failures,
		Success:  len(durations),
	}

	if len(durations) > 0 {
		s.Min = durations[0]
		s.Max = durations[len(durations)-1]
		s.Mean = total / time.Duration(len(durations))
		s.P50 = percentile(durations, 50)
		s.P95 = percentile(durations, 95)
		s.P99 = percentile(durations, 99)
	}

	return s
}

func percentile(sorted []time.Duration, p int) time.Duration {
	idx := int(float64(len(sorted)-1) * float64(p) / 100.0)
	return sorted[idx]
}

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Summary) {
	fmt.Fprintf(w, "\n--- Results ---\n")
	fmt.Fprintf(w, "Total:    %d\n", s.Total)
	fmt.Fprintf(w, "Success:  %d\n", s.Success)
	fmt.Fprintf(w, "Failures: %d\n", s.Failures)
	fmt.Fprintf(w, "Mean:     %v\n", s.Mean)
	fmt.Fprintf(w, "Min:      %v\n", s.Min)
	fmt.Fprintf(w, "Max:      %v\n", s.Max)
	fmt.Fprintf(w, "p50:      %v\n", s.P50)
	fmt.Fprintf(w, "p95:      %v\n", s.P95)
	fmt.Fprintf(w, "p99:      %v\n", s.P99)
}

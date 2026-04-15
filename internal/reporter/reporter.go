package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/bojand/grpcannon/internal/stats"
)

// Format defines the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter handles formatting and writing load test results.
type Reporter struct {
	format Format
	out    io.Writer
}

// New creates a new Reporter with the given format and output writer.
// If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{format: format, out: out}
}

// Report writes the computed stats summary to the output.
func (r *Reporter) Report(summary stats.Summary, elapsed time.Duration) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(summary, elapsed)
	default:
		return r.writeText(summary, elapsed)
	}
}

func (r *Reporter) writeText(s stats.Summary, elapsed time.Duration) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n--- gRPCannon Report ---\n")
	fmt.Fprintf(w, "Duration:\t%s\n", elapsed.Round(time.Millisecond))
	fmt.Fprintf(w, "Total:\t%d\n", s.Total)
	fmt.Fprintf(w, "Succeeded:\t%d\n", s.Succeeded)
	fmt.Fprintf(w, "Failed:\t%d\n", s.Failed)
	fmt.Fprintf(w, "RPS:\t%.2f\n", float64(s.Total)/elapsed.Seconds())
	fmt.Fprintf(w, "\nLatency:\n")
	fmt.Fprintf(w, "  Min:\t%s\n", s.Min.Round(time.Microsecond))
	fmt.Fprintf(w, "  Mean:\t%s\n", s.Mean.Round(time.Microsecond))
	fmt.Fprintf(w, "  P50:\t%s\n", s.P50.Round(time.Microsecond))
	fmt.Fprintf(w, "  P95:\t%s\n", s.P95.Round(time.Microsecond))
	fmt.Fprintf(w, "  P99:\t%s\n", s.P99.Round(time.Microsecond))
	fmt.Fprintf(w, "  Max:\t%s\n", s.Max.Round(time.Microsecond))
	return w.Flush()
}

func (r *Reporter) writeJSON(s stats.Summary, elapsed time.Duration) error {
	fmt.Fprintf(r.out, `{"duration_ms":%d,"total":%d,"succeeded":%d,"failed":%d,"rps":%.2f,"latency":{"min_us":%d,"mean_us":%d,"p50_us":%d,"p95_us":%d,"p99_us":%d,"max_us":%d}}\n`,
		elapsed.Milliseconds(),
		s.Total, s.Succeeded, s.Failed,
		float64(s.Total)/elapsed.Seconds(),
		s.Min.Microseconds(), s.Mean.Microseconds(),
		s.P50.Microseconds(), s.P95.Microseconds(),
		s.P99.Microseconds(), s.Max.Microseconds(),
	)
	return nil
}

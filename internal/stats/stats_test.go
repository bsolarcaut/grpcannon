package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/grpcannon/internal/runner"
	"github.com/yourusername/grpcannon/internal/stats"
)

func makeResults(durations []time.Duration, errCount int) []runner.Result {
	results := make([]runner.Result, 0, len(durations)+errCount)
	for _, d := range durations {
		results = append(results, runner.Result{Duration: d})
	}
	for i := 0; i < errCount; i++ {
		results = append(results, runner.Result{Err: errSentinel})
	}
	return results
}

var errSentinel = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }

func TestCompute_Empty(t *testing.T) {
	s := stats.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected Total 0, got %d", s.Total)
	}
}

func TestCompute_AllSuccess(t *testing.T) {
	durations := []time.Duration{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	results := makeResults(durations, 0)
	s := stats.Compute(results)

	if s.Total != 10 || s.Success != 10 || s.Failures != 0 {
		t.Errorf("unexpected counts: %+v", s)
	}
	if s.Min != 10 {
		t.Errorf("expected min 10, got %v", s.Min)
	}
	if s.Max != 100 {
		t.Errorf("expected max 100, got %v", s.Max)
	}
	if s.Mean != 55 {
		t.Errorf("expected mean 55, got %v", s.Mean)
	}
}

func TestCompute_WithErrors(t *testing.T) {
	durations := []time.Duration{10, 20, 30}
	results := makeResults(durations, 2)
	s := stats.Compute(results)

	if s.Total != 5 {
		t.Errorf("expected Total 5, got %d", s.Total)
	}
	if s.Failures != 2 {
		t.Errorf("expected Failures 2, got %d", s.Failures)
	}
	if s.Success != 3 {
		t.Errorf("expected Success 3, got %d", s.Success)
	}
}

func TestPrint_ContainsFields(t *testing.T) {
	durations := []time.Duration{
		1 * time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond,
	}
	s := stats.Compute(makeResults(durations, 1))

	var buf bytes.Buffer
	stats.Print(&buf, s)
	out := buf.String()

	for _, field := range []string{"Total", "Success", "Failures", "Mean", "Min", "Max", "p50", "p95", "p99"} {
		if !strings.Contains(out, field) {
			t.Errorf("output missing field %q", field)
		}
	}
}

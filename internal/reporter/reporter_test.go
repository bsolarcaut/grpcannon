package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bojand/grpcannon/internal/reporter"
	"github.com/bojand/grpcannon/internal/stats"
)

func makeSummary() stats.Summary {
	return stats.Summary{
		Total:     100,
		Succeeded: 95,
		Failed:    5,
		Min:       1 * time.Millisecond,
		Mean:      5 * time.Millisecond,
		P50:       4 * time.Millisecond,
		P95:       9 * time.Millisecond,
		P99:       10 * time.Millisecond,
		Max:       12 * time.Millisecond,
	}
}

func TestReporter_Text(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)

	err := r.Report(makeSummary(), 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"gRPCannon Report", "Total:", "Succeeded:", "Failed:", "RPS:", "P95:", "P99:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestReporter_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)

	err := r.Report(makeSummary(), 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{`"total":100`, `"succeeded":95`, `"failed":5`, `"p95_us"`, `"rps"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected JSON to contain %q, got:\n%s", want, out)
		}
	}
}

func TestReporter_DefaultsToStdout(t *testing.T) {
	// Ensure New does not panic with nil writer; just verify no error.
	r := reporter.New(reporter.FormatText, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestReporter_UnknownFormatFallsBackToText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.Format("xml"), &buf)

	err := r.Report(makeSummary(), 1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "gRPCannon Report") {
		t.Error("expected fallback to text format")
	}
}

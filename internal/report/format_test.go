package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/report"
	"github.com/example/grpcannon/internal/stats"
)

func buildReport() report.RunReport {
	b := report.New("localhost:50051", "/svc/Call", 4)
	for i := 0; i < 3; i++ {
		b.Add(stats.Result{Duration: 10 * time.Millisecond})
	}
	return b.Build()
}

func TestPrint_TextContainsTarget(t *testing.T) {
	var buf bytes.Buffer
	if err := report.Print(&buf, buildReport(), report.FormatText); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "localhost:50051") {
		t.Fatalf("text output missing target: %s", buf.String())
	}
}

func TestPrint_TextContainsRPS(t *testing.T) {
	var buf bytes.Buffer
	report.Print(&buf, buildReport(), report.FormatText)
	if !strings.Contains(buf.String(), "RPS") {
		t.Fatalf("text output missing RPS field")
	}
}

func TestPrint_JSONIsValid(t *testing.T) {
	var buf bytes.Buffer
	if err := report.Print(&buf, buildReport(), report.FormatJSON); err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrint_JSONContainsMethod(t *testing.T) {
	var buf bytes.Buffer
	report.Print(&buf, buildReport(), report.FormatJSON)
	if !strings.Contains(buf.String(), "/svc/Call") {
		t.Fatalf("JSON missing method field")
	}
}

func TestPrint_UnknownFormatFallsBackToText(t *testing.T) {
	var buf bytes.Buffer
	if err := report.Print(&buf, buildReport(), report.Format("xml")); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Target") {
		t.Fatalf("fallback text missing Target field")
	}
}

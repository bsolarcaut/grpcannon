package report_test

import (
	"errors"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/report"
	"github.com/example/grpcannon/internal/stats"
)

func makeResult(dur time.Duration, err error) stats.Result {
	return stats.Result{Duration: dur, Err: err}
}

func TestBuild_EmptyResults(t *testing.T) {
	b := report.New("localhost:50051", "/pkg.Svc/Method", 4)
	r := b.Build()
	if r.Summary.Total != 0 {
		t.Fatalf("expected 0 total, got %d", r.Summary.Total)
	}
	if r.RPS != 0 {
		t.Fatalf("expected 0 rps, got %f", r.RPS)
	}
}

func TestBuild_MetadataPreserved(t *testing.T) {
	b := report.New("host:9090", "/svc/Method", 8)
	r := b.Build()
	if r.Target != "host:9090" {
		t.Fatalf("unexpected target: %s", r.Target)
	}
	if r.Method != "/svc/Method" {
		t.Fatalf("unexpected method: %s", r.Method)
	}
	if r.Workers != 8 {
		t.Fatalf("unexpected workers: %d", r.Workers)
	}
}

func TestBuild_TotalMatchesAdded(t *testing.T) {
	b := report.New("h:1", "/m", 1)
	for i := 0; i < 10; i++ {
		b.Add(makeResult(10*time.Millisecond, nil))
	}
	r := b.Build()
	if r.Summary.Total != 10 {
		t.Fatalf("expected 10 total, got %d", r.Summary.Total)
	}
}

func TestBuild_ErrorsReflectedInSummary(t *testing.T) {
	b := report.New("h:1", "/m", 2)
	b.Add(makeResult(5*time.Millisecond, nil))
	b.Add(makeResult(5*time.Millisecond, errors.New("boom")))
	r := b.Build()
	if r.Summary.Errors != 1 {
		t.Fatalf("expected 1 error, got %d", r.Summary.Errors)
	}
}

func TestBuild_RPSPositiveWhenResultsExist(t *testing.T) {
	b := report.New("h:1", "/m", 1)
	for i := 0; i < 5; i++ {
		b.Add(makeResult(1*time.Millisecond, nil))
	}
	r := b.Build()
	if r.RPS <= 0 {
		t.Fatalf("expected positive RPS, got %f", r.RPS)
	}
}

func TestBuild_DurationPositive(t *testing.T) {
	b := report.New("h:1", "/m", 1)
	b.Add(makeResult(1*time.Millisecond, nil))
	r := b.Build()
	if r.Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", r.Duration)
	}
}

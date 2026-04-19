package scoreboard_test

import (
	"testing"
	"time"

	"github.com/example/grpcannon/internal/scoreboard"
)

func TestNew_Empty(t *testing.T) {
	b := scoreboard.New()
	if got := b.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestRecord_SingleWorker(t *testing.T) {
	b := scoreboard.New()
	b.Record(1, 10*time.Millisecond, false)
	b.Record(1, 20*time.Millisecond, false)

	snap := b.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	e := snap[0]
	if e.Total != 2 {
		t.Errorf("expected Total=2, got %d", e.Total)
	}
	if e.Errors != 0 {
		t.Errorf("expected Errors=0, got %d", e.Errors)
	}
	if e.AvgLatency() != 15*time.Millisecond {
		t.Errorf("expected AvgLatency=15ms, got %v", e.AvgLatency())
	}
}

func TestRecord_ErrorCounted(t *testing.T) {
	b := scoreboard.New()
	b.Record(2, 5*time.Millisecond, true)
	b.Record(2, 5*time.Millisecond, false)

	snap := b.Snapshot()
	if snap[0].Errors != 1 {
		t.Errorf("expected 1 error, got %d", snap[0].Errors)
	}
}

func TestRecord_MultipleWorkers(t *testing.T) {
	b := scoreboard.New()
	for i := 0; i < 5; i++ {
		b.Record(i, time.Millisecond, false)
	}
	if got := len(b.Snapshot()); got != 5 {
		t.Errorf("expected 5 workers, got %d", got)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	b := scoreboard.New()
	b.Record(1, time.Millisecond, false)
	b.Reset()
	if got := len(b.Snapshot()); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}

func TestAvgLatency_ZeroTotal(t *testing.T) {
	e := scoreboard.Entry{}
	if e.AvgLatency() != 0 {
		t.Errorf("expected 0 latency for empty entry")
	}
}

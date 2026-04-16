package snapshot_test

import (
	"sync"
	"testing"

	"github.com/nickcorin/grpcannon/internal/snapshot"
)

func TestSnap_InitialState(t *testing.T) {
	c := snapshot.NewCollector()
	m := c.Snap()
	if m.Total != 0 || m.Successes != 0 || m.Errors != 0 {
		t.Fatalf("expected zero counters, got %+v", m)
	}
}

func TestRecordSuccess(t *testing.T) {
	c := snapshot.NewCollector()
	c.RecordSuccess()
	c.RecordSuccess()
	m := c.Snap()
	if m.Total != 2 || m.Successes != 2 || m.Errors != 0 {
		t.Fatalf("unexpected metrics: %+v", m)
	}
}

func TestRecordError(t *testing.T) {
	c := snapshot.NewCollector()
	c.RecordError()
	m := c.Snap()
	if m.Total != 1 || m.Errors != 1 || m.Successes != 0 {
		t.Fatalf("unexpected metrics: %+v", m)
	}
}

func TestSetInFlight(t *testing.T) {
	c := snapshot.NewCollector()
	c.SetInFlight(7)
	if got := c.Snap().InFlight; got != 7 {
		t.Fatalf("expected 7 in-flight, got %d", got)
	}
}

func TestSnap_Elapsed(t *testing.T) {
	c := snapshot.NewCollector()
	m := c.Snap()
	if m.Elapsed < 0 {
		t.Fatal("elapsed should be non-negative")
	}
}

func TestCollector_ConcurrentSafety(t *testing.T) {
	c := snapshot.NewCollector()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); c.RecordSuccess() }()
		go func() { defer wg.Done(); c.RecordError() }()
	}
	wg.Wait()
	m := c.Snap()
	if m.Total != 200 {
		t.Fatalf("expected 200 total, got %d", m.Total)
	}
	if m.Successes+m.Errors != m.Total {
		t.Fatalf("successes+errors != total: %+v", m)
	}
}

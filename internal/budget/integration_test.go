package budget_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/budget"
)

func TestIntegration_ConcurrentRecordAndCheck(t *testing.T) {
	b := budget.New(0.50, 20)
	var wg sync.WaitGroup
	const workers = 8
	const each = 50

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < each; i++ {
				if i%4 == 0 {
					b.Record(errors.New("err"))
				} else {
					b.Record(nil)
				}
			}
		}(w)
	}
	wg.Wait()

	total, failures := b.Snapshot()
	if total != workers*each {
		t.Fatalf("expected %d total, got %d", workers*each, total)
	}
	expectedFailures := int64(workers * each / 4)
	if failures != expectedFailures {
		t.Fatalf("expected %d failures, got %d", expectedFailures, failures)
	}
	// 25% failure rate < 50% threshold → should be OK
	if err := b.Check(); err != nil {
		t.Fatalf("unexpected exhaustion: %v", err)
	}
}

// TestIntegration_BudgetExhausted verifies that Check returns an error when the
// failure rate exceeds the configured threshold.
func TestIntegration_BudgetExhausted(t *testing.T) {
	b := budget.New(0.25, 10)

	// Record 3 failures out of 4 total → 75% failure rate > 25% threshold.
	b.Record(errors.New("err"))
	b.Record(errors.New("err"))
	b.Record(errors.New("err"))
	b.Record(nil)

	if err := b.Check(); err == nil {
		t.Fatal("expected budget exhaustion error, got nil")
	}
}

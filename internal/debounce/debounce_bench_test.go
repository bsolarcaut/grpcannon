package debounce_test

import (
	"testing"
	"time"

	"github.com/example/grpcannon/internal/debounce"
)

// BenchmarkTrigger measures the overhead of repeated Trigger calls that keep
// resetting the timer without ever firing.
func BenchmarkTrigger(b *testing.B) {
	d := debounce.New(10*time.Second, func() {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Trigger()
	}
	d.Stop()
}

// BenchmarkFlush measures Trigger followed immediately by Flush.
func BenchmarkFlush(b *testing.B) {
	var calls int
	d := debounce.New(10*time.Second, func() { calls++ })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Trigger()
		d.Flush()
	}
	_ = calls
}

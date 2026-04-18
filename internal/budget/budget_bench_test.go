package budget_test

import (
	"testing"

	"github.com/example/grpcannon/internal/budget"
)

func BenchmarkRecord(b *testing.B) {
	bg := budget.New(0.05, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bg.Record(nil)
	}
}

func BenchmarkCheck(b *testing.B) {
	bg := budget.New(0.05, 1)
	bg.Record(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bg.Check()
	}
}

func BenchmarkRecord_Parallel(b *testing.B) {
	bg := budget.New(0.05, 10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bg.Record(nil)
		}
	})
}

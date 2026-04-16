package window

import (
	"testing"
	"time"
)

func BenchmarkAdd(b *testing.B) {
	w := New(60, time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Add(1)
	}
}

func BenchmarkSum(b *testing.B) {
	w := New(60, time.Second)
	for i := 0; i < 1000; i++ {
		w.Add(1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Sum()
	}
}

func BenchmarkRate(b *testing.B) {
	w := New(60, time.Second)
	for i := 0; i < 1000; i++ {
		w.Add(1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Rate()
	}
}

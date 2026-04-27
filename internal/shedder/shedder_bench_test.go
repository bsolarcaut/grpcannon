package shedder_test

import (
	"testing"

	"github.com/yourorg/grpcannon/internal/shedder"
)

func BenchmarkAllow(b *testing.B) {
	s := shedder.New(shedder.Default())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.Allow()
		}
	})
}

func BenchmarkRecord(b *testing.B) {
	s := shedder.New(shedder.Default())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Record(false)
		}
	})
}

func BenchmarkAllowRecord(b *testing.B) {
	s := shedder.New(shedder.Default())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Allow()
		s.Record(i%10 == 0)
	}
}

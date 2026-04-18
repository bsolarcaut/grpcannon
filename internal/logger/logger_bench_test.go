package logger

import (
	"io"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	l := New(io.Discard, LevelInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("request completed in %dms", 12)
	}
}

func BenchmarkDebugSuppressed(b *testing.B) {
	l := New(io.Discard, LevelInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("verbose detail %d", i)
	}
}

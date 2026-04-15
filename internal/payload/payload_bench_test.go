package payload_test

import (
	"testing"

	"github.com/user/grpcannon/internal/payload"
)

// BenchmarkFromJSON measures the cost of parsing a moderately sized JSON object.
func BenchmarkFromJSON(b *testing.B) {
	const raw = `{"name":"benchmark","iterations":1000,"enabled":true,"ratio":0.75}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := payload.FromJSON(raw)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFromPairs measures the cost of parsing key=value pairs.
func BenchmarkFromPairs(b *testing.B) {
	pairs := []string{"name=benchmark", "iterations=1000", "enabled=true", "ratio=0.75"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := payload.New()
		if err := builder.FromPairs(pairs); err != nil {
			b.Fatal(err)
		}
		_ = builder.Build()
n}

// BenchmarkJSON measures JSON serialisation of a built payload.
func BenchmarkJSON(b *testing.B) {
	builder := payload.New()
	_ = builder.FromPairs([]string{"a=1", "b=2", "c=3"})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := builder.JSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}

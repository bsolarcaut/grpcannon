package drain_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/drain"
)

func BenchmarkAcquireRelease(b *testing.B) {
	d := drain.New(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d.Acquire() {
			d.Release()
		}
	}
	_ = d.Drain(context.Background())
}

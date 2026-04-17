package probe_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/grpcannon/internal/probe"
)

func BenchmarkCheck(b *testing.B) {
	conn := localConn(b.(*testing.B))
	p, err := probe.New(conn, "/svc/Method", 100*time.Millisecond)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Check(ctx)
	}
}

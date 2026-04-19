package barrier_test

import (
	"context"
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/barrier"
)

// BenchmarkWait measures the throughput of a barrier with GOMAXPROCS workers
// cycling through many generations.
func BenchmarkWait(b *testing.B) {
	const workers = 8
	br := barrier.New(workers)
	var wg sync.WaitGroup
	start := make(chan struct{})

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < b.N; i++ {
				_ = br.Wait(context.Background())
			}
		}()
	}

	b.ResetTimer()
	close(start)
	wg.Wait()
}

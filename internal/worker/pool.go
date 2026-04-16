package worker

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// Pool manages a fixed number of concurrent Workers.
type Pool struct {
	concurrency int
	conn        *grpc.ClientConn
	callFn      CallFunc
}

// NewPool creates a Pool with the given concurrency level.
func NewPool(concurrency int, conn *grpc.ClientConn, fn CallFunc) *Pool {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Pool{
		concurrency: concurrency,
		conn:        conn,
		callFn:      fn,
	}
}

// Run spawns concurrency workers, distributes totalRequests among them, and
// collects all Results into the returned slice. It blocks until all workers
// finish or the context is cancelled.
func (p *Pool) Run(ctx context.Context, totalRequests int) []Result {
	resultsCh := make(chan Result, totalRequests)

	perWorker := totalRequests / p.concurrency
	remainder := totalRequests % p.concurrency

	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		n := perWorker
		if i == 0 {
			n += remainder
		}
		if n == 0 {
			continue
		}
		wg.Add(1)
		go func(requests int) {
			defer wg.Done()
			w := New(p.conn, p.callFn, resultsCh)
			w.Run(ctx, requests)
		}(n)
	}

	wg.Wait()
	close(resultsCh)

	results := make([]Result, 0, totalRequests)
	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}

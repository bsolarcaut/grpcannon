// Package warmup implements a pre-benchmark warm-up phase for grpcannon.
//
// A warm-up phase primes the target service and any intermediate proxies
// before latency measurements begin, reducing the impact of cold-start
// effects on histogram data.
//
// Usage:
//
//	cfg := warmup.DefaultConfig()
//	cfg.Requests = 20
//	if err := warmup.Run(ctx, invoker, cfg); err != nil {
//		log.Fatalf("warm-up failed: %v", err)
//	}
package warmup

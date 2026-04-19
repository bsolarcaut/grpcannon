// Package observe implements a periodic metrics observer for grpcannon load
// runs. It reads live counters from a [snapshot.Collector] and writes
// human-readable progress lines at a configurable interval so operators can
// watch throughput and error rates without waiting for the final report.
package observe

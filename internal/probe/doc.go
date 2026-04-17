// Package probe implements a pre-flight health check for gRPC targets.
//
// Usage:
//
//	p, err := probe.New(conn, "/pkg.Service/Method", 3*time.Second)
//	if err != nil { ... }
//	result := p.CheckN(ctx, 3, 500*time.Millisecond)
//	if !result.OK { log.Fatal(result.Err) }
//
// CheckN retries up to n times with a configurable interval, returning
// immediately on the first successful response.
package probe

// Package circuit provides a thread-safe circuit breaker used to protect
// gRPC load-test workers from cascading failures. When the number of
// consecutive errors exceeds a configurable threshold the breaker opens and
// subsequent calls are rejected immediately with ErrOpen. After a cooldown
// period the breaker transitions to half-open and allows a single probe
// request through; a successful response closes the breaker again.
package circuit

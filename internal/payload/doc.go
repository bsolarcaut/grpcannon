// Package payload provides helpers for constructing gRPC request
// payloads used by grpcannon during load tests.
//
// A payload can be sourced from:
//   - a raw JSON string (--data flag)
//   - a list of key=value pairs (--data-pair flags)
//
// Both sources can be combined; key=value pairs take precedence over
// JSON fields with the same key when FromPairs is called after FromJSON.
//
// Example usage:
//
//	b, err := payload.FromJSON(rawJSON)
//	if err != nil { ... }
//	if err := b.FromPairs(pairs); err != nil { ... }
//	data := b.Build()
package payload

// Package timeout provides small utilities for managing per-call deadlines
// and duration defaults used throughout grpcannon.
//
// # WithDeadline
//
// WithDeadline wraps context.WithTimeout and validates that the supplied
// duration is positive before creating the child context.
//
// # Clamp
//
// Clamp constrains a duration to a [min, max] range, useful for bounding
// user-supplied timeout values before passing them to gRPC calls.
//
// # Default
//
// Default selects a fallback duration when the caller supplies zero,
// mirroring the pattern used in grpcclient and invoker packages.
package timeout

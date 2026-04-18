// Package logger provides a minimal levelled logger used throughout grpcannon.
//
// Usage:
//
//	l := logger.New(os.Stderr, logger.LevelInfo)
//	l.Info("starting run with %d workers", concurrency)
//
// Levels (ascending severity): Debug, Info, Warn, Error.
// Messages below the configured minimum level are silently dropped.
package logger

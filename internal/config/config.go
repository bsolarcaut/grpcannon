package config

import (
	"errors"
	"time"
)

// Config holds all runtime configuration for a grpcannon load test.
type Config struct {
	// Target is the gRPC server address (host:port).
	Target string

	// Method is the fully-qualified gRPC method name, e.g. "package.Service/Method".
	Method string

	// Concurrency is the number of concurrent workers sending requests.
	Concurrency int

	// TotalRequests is the total number of requests to send across all workers.
	TotalRequests int

	// Duration overrides TotalRequests when non-zero; the test runs for this long.
	Duration time.Duration

	// Timeout is the per-request deadline.
	Timeout time.Duration

	// Insecure disables TLS verification when true.
	Insecure bool

	// Metadata holds key=value pairs forwarded as gRPC metadata.
	Metadata map[string]string

	// PayloadJSON is the JSON-encoded request payload.
	PayloadJSON string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Concurrency:   10,
		TotalRequests: 200,
		Timeout:       5 * time.Second,
		Insecure:      false,
		Metadata:      make(map[string]string),
	}
}

// Validate checks that the configuration is self-consistent.
func (c *Config) Validate() error {
	if c.Target == "" {
		return errors.New("target address must not be empty")
	}
	if c.Method == "" {
		return errors.New("gRPC method must not be empty")
	}
	if c.Concurrency <= 0 {
		return errors.New("concurrency must be greater than 0")
	}
	if c.Duration == 0 && c.TotalRequests <= 0 {
		return errors.New("either duration or total-requests must be specified")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	return nil
}

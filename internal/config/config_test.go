package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Concurrency != 10 {
		t.Errorf("expected default concurrency 10, got %d", cfg.Concurrency)
	}
	if cfg.TotalRequests != 200 {
		t.Errorf("expected default total requests 200, got %d", cfg.TotalRequests)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Timeout)
	}
	if cfg.Insecure {
		t.Error("expected insecure to be false by default")
	}
	if cfg.Metadata == nil {
		t.Error("expected metadata map to be initialised")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost:50051"
	cfg.Method = "helloworld.Greeter/SayHello"

	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Method = "helloworld.Greeter/SayHello"

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing target")
	}
}

func TestValidate_MissingMethod(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost:50051"

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing method")
	}
}

func TestValidate_InvalidConcurrency(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost:50051"
	cfg.Method = "svc/Method"
	cfg.Concurrency = 0

	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero concurrency")
	}
}

func TestValidate_DurationOverridesTotalRequests(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "localhost:50051"
	cfg.Method = "svc/Method"
	cfg.TotalRequests = 0
	cfg.Duration = 10 * time.Second

	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error when duration is set: %v", err)
	}
}

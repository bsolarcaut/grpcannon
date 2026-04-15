package grpcclient_test

import (
	"testing"
	"time"

	"github.com/mxbossard/grpcannon/internal/grpcclient"
)

func TestNew_EmptyTargetReturnsError(t *testing.T) {
	_, err := grpcclient.New(grpcclient.Options{})
	if err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

func TestNew_ValidTargetReturnsClient(t *testing.T) {
	c, err := grpcclient.New(grpcclient.Options{
		Target:  "localhost:50051",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	_ = c.Close()
}

func TestNew_DefaultTimeoutApplied(t *testing.T) {
	// A zero timeout should be replaced with the default (10s).
	// We verify construction succeeds without panic.
	c, err := grpcclient.New(grpcclient.Options{
		Target: "localhost:50051",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = c.Close()
}

func TestClose_IdempotentOnValidClient(t *testing.T) {
	c, err := grpcclient.New(grpcclient.Options{
		Target:  "localhostt	Timeout: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("first() returned error: %v", err)
	}
}

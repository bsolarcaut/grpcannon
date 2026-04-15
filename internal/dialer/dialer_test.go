package dialer_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/nickcorin/grpcannon/internal/dialer"
	"google.golang.org/grpc"
)

func startLocalServer(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestDial_EmptyTargetReturnsError(t *testing.T) {
	_, err := dialer.Dial(context.Background(), dialer.Options{})
	if err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

func TestDial_ValidTargetConnects(t *testing.T) {
	addr := startLocalServer(t)
	conn, err := dialer.Dial(context.Background(), dialer.Options{
		Target:      addr,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer conn.Close()

	if conn.ClientConn == nil {
		t.Fatal("expected non-nil ClientConn")
	}
}

func TestDial_TimeoutOnUnreachableTarget(t *testing.T) {
	ctx := context.Background()
	_, err := dialer.Dial(ctx, dialer.Options{
		Target:      "127.0.0.1:19999",
		DialTimeout: 300 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestDial_DefaultTimeoutApplied(t *testing.T) {
	addr := startLocalServer(t)
	// DialTimeout zero should use default (10s), connection should succeed.
	conn, err := dialer.Dial(context.Background(), dialer.Options{
		Target: addr,
	})
	if err != nil {
		t.Fatalf("unexpected error with default timeout: %v", err)
	}
	defer conn.Close()
}

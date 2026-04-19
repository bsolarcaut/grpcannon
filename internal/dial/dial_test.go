package dial_test

import (
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/example/grpcannon/internal/dial"
)

func startServer(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(func() { srv.Stop() })
	return lis.Addr().String()
}

func TestOpen_EmptyTargetReturnsError(t *testing.T) {
	_, err := dial.Open(dial.Options{})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestOpen_ValidTargetReturnsConn(t *testing.T) {
	addr := startServer(t)
	c, err := dial.Open(dial.Options{Target: addr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer c.Close()
	if c.Target != addr {
		t.Errorf("target mismatch: got %q want %q", c.Target, addr)
	}
}

func TestOpen_DefaultTimeoutApplied(t *testing.T) {
	addr := startServer(t)
	c, err := dial.Open(dial.Options{Target: addr, Timeout: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer c.Close()
}

func TestOpen_TimeoutOnUnreachable(t *testing.T) {
	_, err := dial.Open(dial.Options{
		Target:  "127.0.0.1:1",
		Timeout: 100 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error for unreachable target")
	}
}

func TestClose_NilConnIsNoop(t *testing.T) {
	var c *dial.Conn
	if err := c.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

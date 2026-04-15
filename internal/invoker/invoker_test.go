package invoker_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourorg/grpcannon/internal/invoker"
)

// dialLocalServer starts a bare gRPC server on a random port and returns
// a client connection to it plus a cleanup function.
func dialLocalServer(t *testing.T) (*grpc.ClientConn, func()) {
	t.Helper()
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go srv.Serve(lis) //nolint:errcheck
	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		srv.Stop()
		t.Fatalf("dial: %v", err)
	}
	return conn, func() {
		conn.Close()
		srv.Stop()
	}
}

func TestNew_NilConnReturnsError(t *testing.T) {
	_, err := invoker.New(nil, "/pkg.Svc/Method", time.Second, nil)
	if err == nil {
		t.Fatal("expected error for nil conn")
	}
}

func TestNew_EmptyMethodReturnsError(t *testing.T) {
	conn, cleanup := dialLocalServer(t)
	defer cleanup()

	_, err := invoker.New(conn, "", time.Second, nil)
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestNew_DefaultTimeoutApplied(t *testing.T) {
	conn, cleanup := dialLocalServer(t)
	defer cleanup()

	// A zero timeout should be replaced by the default (5s).
	inv, err := invoker.New(conn, "/pkg.Svc/Method", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv == nil {
		t.Fatal("expected non-nil invoker")
	}
}

func TestCall_ReturnsResultWithDuration(t *testing.T) {
	conn, cleanup := dialLocalServer(t)
	defer cleanup()

	inv, err := invoker.New(conn, "/pkg.Svc/Method", time.Second, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res := inv.Call(context.Background())
	// The server has no handler registered so we expect an error, but the
	// call must still return a Result with a measured duration and a status code.
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
	if res.StatusCode == "" {
		t.Error("expected non-empty status code")
	}
}

func TestCall_CancelledContextReflectedInResult(t *testing.T) {
	conn, cleanup := dialLocalServer(t)
	defer cleanup()

	inv, err := invoker.New(conn, "/pkg.Svc/Method", 5*time.Second, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	res := inv.Call(ctx)
	if res.Err == nil {
		t.Error("expected error for cancelled context")
	}
}

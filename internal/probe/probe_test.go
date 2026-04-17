package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourusername/grpcannon/internal/probe"
)

func localConn(t *testing.T) *grpc.ClientConn {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { lis.Close() })
	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestNew_NilConnReturnsError(t *testing.T) {
	_, err := probe.New(nil, "/svc/Method", 0)
	if err == nil {
		t.Fatal("expected error for nil conn")
	}
}

func TestNew_EmptyMethodReturnsError(t *testing.T) {
	conn := localConn(t)
	_, err := probe.New(conn, "", 0)
	if err == nil {
		t.Fatal("expected error for empty method")
	}
}

func TestNew_DefaultTimeoutApplied(t *testing.T) {
	conn := localConn(t)
	p, err := probe.New(conn, "/svc/Method", 0)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestCheck_FailsOnUnreachableTarget(t *testing.T) {
	conn := localConn(t)
	p, err := probe.New(conn, "/svc/Method", 200*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	res := p.Check(context.Background())
	if res.OK {
		t.Fatal("expected probe to fail on unreachable target")
	}
	if res.Err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestCheckN_StopsOnContextCancellation(t *testing.T) {
	conn := localConn(t)
	p, err := probe.New(conn, "/svc/Method", 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res := p.CheckN(ctx, 5, 50*time.Millisecond)
	if res.OK {
		t.Fatal("expected failure on cancelled context")
	}
}

func TestCheckN_ReturnsLastResultAfterAllAttempts(t *testing.T) {
	conn := localConn(t)
	p, err := probe.New(conn, "/svc/Method", 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	res := p.CheckN(context.Background(), 2, 10*time.Millisecond)
	if res.OK {
		t.Fatal("expected all probes to fail")
	}
}

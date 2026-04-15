// Package invoker provides a gRPC method invoker that uses dynamic reflection
// to call arbitrary gRPC methods without generated stubs.
package invoker

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Result holds the outcome of a single gRPC invocation.
type Result struct {
	// Duration is the round-trip latency of the call.
	Duration time.Duration
	// Err is non-nil if the call failed.
	Err error
	// StatusCode is the gRPC status code string (e.g. "OK", "UNAVAILABLE").
	StatusCode string
}

// Invoker executes a single gRPC unary call against a pre-dialed connection.
type Invoker struct {
	conn    *grpc.ClientConn
	method  string
	timeout time.Duration
	headers map[string]string
}

// New returns an Invoker configured for the given connection and method.
// method must be the fully-qualified gRPC path, e.g. "/pkg.Service/Method".
func New(conn *grpc.ClientConn, method string, timeout time.Duration, headers map[string]string) (*Invoker, error) {
	if conn == nil {
		return nil, fmt.Errorf("invoker: conn must not be nil")
	}
	if method == "" {
		return nil, fmt.Errorf("invoker: method must not be empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Invoker{
		conn:    conn,
		method:  method,
		timeout: timeout,
		headers: headers,
	}, nil
}

// Call performs a single unary gRPC invocation and returns a Result.
// The caller-supplied context is respected for cancellation; an additional
// per-call deadline is applied on top via the configured timeout.
func (inv *Invoker) Call(ctx context.Context) Result {
	callCtx, cancel := context.WithTimeout(ctx, inv.timeout)
	defer cancel()

	if len(inv.headers) > 0 {
		md := metadata.New(inv.headers)
		callCtx = metadata.NewOutgoingContext(callCtx, md)
	}

	start := time.Now()
	err := inv.conn.Invoke(callCtx, inv.method, struct{}{}, &struct{}{})
	dur := time.Since(start)

	code := grpcStatusCode(err)
	return Result{
		Duration:   dur,
		Err:        err,
		StatusCode: code,
	}
}

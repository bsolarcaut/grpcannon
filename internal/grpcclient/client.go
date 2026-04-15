package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// CallFunc is the signature for a gRPC unary call used by the runner.
type CallFunc func(ctx context.Context) error

// Options holds configuration for the gRPC client connection.
type Options struct {
	Target  string
	Timeout time.Duration
	Headers map[string]string
}

// Client wraps a gRPC client connection.
type Client struct {
	conn *grpc.ClientConn
	opts Options
}

// New dials the target and returns a connected Client.
func New(opts Options) (*Client, error) {
	if opts.Target == "" {
		return nil, fmt.Errorf("grpcclient: target must not be empty")
	}
	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	conn, err := grpc.NewClient(
		opts.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpcclient: dial %s: %w", opts.Target, err)
	}

	return &Client{conn: conn, opts: opts}, nil
}

// Invoke performs a raw unary RPC call on the given fully-qualified method.
// The method must be in the format "/package.Service/Method".
func (c *Client) Invoke(ctx context.Context, method string, req, reply interface{}) error {
	if len(c.opts.Headers) > 0 {
		md := metadata.New(c.opts.Headers)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	callCtx, cancel := context.WithTimeout(ctx, c.opts.Timeout)
	defer cancel()

	return c.conn.Invoke(callCtx, method, req, reply)
}

// Close tears down the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

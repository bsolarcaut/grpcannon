// Package dial provides a reusable connection factory that wraps
// dialer and grpcclient construction behind a single call.
package dial

import (
	"errors"
	"time"

	"google.golang.org/grpc"

	"github.com/example/grpcannon/internal/dialer"
)

// Options configures the connection factory.
type Options struct {
	Target  string
	Timeout time.Duration
	Insecure bool
}

// Conn wraps a *grpc.ClientConn with its originating target.
type Conn struct {
	*grpc.ClientConn
	Target string
}

// Open dials target and returns a Conn ready for use.
func Open(opts Options) (*Conn, error) {
	if opts.Target == "" {
		return nil, errors.New("dial: target must not be empty")
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 5 * time.Second
	}
	cc, err := dialer.Dial(opts.Target, opts.Timeout)
	if err != nil {
		return nil, err
	}
	return &Conn{ClientConn: cc, Target: opts.Target}, nil
}

// Close releases the underlying connection.
func (c *Conn) Close() error {
	if c == nil || c.ClientConn == nil {
		return nil
	}
	return c.ClientConn.Close()
}

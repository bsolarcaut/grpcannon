package dialer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Options holds dialer configuration.
type Options struct {
	Target     string
	TLS        bool
	TLSConfig  *tls.Config
	DialTimeout time.Duration
}

// Conn wraps a gRPC client connection.
type Conn struct {
	*grpc.ClientConn
}

// Dial establishes a gRPC connection using the provided options.
func Dial(ctx context.Context, opts Options) (*Conn, error) {
	if opts.Target == "" {
		return nil, fmt.Errorf("dialer: target must not be empty")
	}

	if opts.DialTimeout <= 0 {
		opts.DialTimeout = 10 * time.Second
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if opts.TLS {
		tlsCfg := opts.TLSConfig
		if tlsCfg == nil {
			tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(ctx, opts.DialTimeout)
	defer cancel()

	cc, err := grpc.DialContext(ctx, opts.Target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dialer: failed to connect to %q: %w", opts.Target, err)
	}

	return &Conn{ClientConn: cc}, nil
}

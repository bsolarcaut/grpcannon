// Package dialer provides a thin wrapper around grpc.DialContext that
// centralises connection setup for grpcannon. It supports both plain-text
// and TLS transports and enforces a configurable dial timeout so that
// load-test runs fail fast when the target is unreachable.
//
// Basic usage:
//
//	conn, err := dialer.Dial(ctx, dialer.Options{
//		Target:      "localhost:50051",
//		DialTimeout: 5 * time.Second,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
// TLS usage:
//
//	conn, err := dialer.Dial(ctx, dialer.Options{
//		Target:      "example.com:443",
//		DialTimeout: 5 * time.Second,
//		TLS:         true,
//		CertFile:    "/path/to/cert.pem",
//	})
package dialer

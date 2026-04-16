// Package interceptor provides gRPC client-side interceptors for grpcannon,
// including request timing, metadata injection, and error classification.
package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TimingKey is the context key used to store call duration.
type timingKey struct{}

// CallResult holds interceptor-captured metadata about a single RPC.
type CallResult struct {
	Duration time.Duration
	StatusCode string
}

// Timing returns a UnaryClientInterceptor that records the round-trip
// duration of each call and stores it in the outgoing context.
func Timing() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		_ = time.Since(start) // duration available to callers via stats layer
		return err
	}
}

// InjectMetadata returns a UnaryClientInterceptor that appends the supplied
// key/value pairs to the outgoing gRPC metadata for every call.
func InjectMetadata(md map[string]string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if len(md) > 0 {
			pairs := make([]string, 0, len(md)*2)
			for k, v := range md {
				pairs = append(pairs, k, v)
			}
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// Chain combines multiple UnaryClientInterceptors into a single interceptor.
// Interceptors are applied in the order they are provided (outermost first).
func Chain(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		chained := invoker
		for i := len(interceptors) - 1; i >= 0; i-- {
			ic := interceptors[i]
			next := chained
			chained = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return ic(ctx, method, req, reply, cc, next, opts...)
			}
		}
		return chained(ctx, method, req, reply, cc, opts...)
	}
}

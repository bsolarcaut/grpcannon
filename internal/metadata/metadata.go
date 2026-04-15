// Package metadata provides utilities for attaching gRPC metadata
// (headers) to outgoing requests during load testing.
package metadata

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

// FromSlice parses a slice of "key:value" strings into a gRPC metadata.MD map.
// Each entry must contain exactly one colon separating the key and value.
// Keys are trimmed and lowercased; values are trimmed.
func FromSlice(pairs []string) (metadata.MD, error) {
	md := metadata.MD{}
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid metadata entry %q: expected \"key:value\" format", pair)
		}
		key := strings.TrimSpace(strings.ToLower(parts[0]))
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("invalid metadata entry %q: key must not be empty", pair)
		}
		md[key] = append(md[key], val)
	}
	return md, nil
}

// Merge combines two metadata.MD maps into a new map. Values from b override
// keys that already exist in a.
func Merge(a, b metadata.MD) metadata.MD {
	out := metadata.MD{}
	for k, v := range a {
		out[k] = append([]string(nil), v...)
	}
	for k, v := range b {
		out[k] = append([]string(nil), v...)
	}
	return out
}

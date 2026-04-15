// Package payload provides utilities for building and validating
// gRPC request payloads from JSON input or flag-supplied key=value pairs.
package payload

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Builder constructs a map[string]interface{} payload that can be
// serialised and forwarded to a gRPC invoker.
type Builder struct {
	fields map[string]interface{}
}

// New returns an empty Builder.
func New() *Builder {
	return &Builder{fields: make(map[string]interface{})}
}

// FromJSON parses a raw JSON object string and returns a Builder
// pre-populated with the decoded fields.
func FromJSON(raw string) (*Builder, error) {
	if strings.TrimSpace(raw) == "" {
		return New(), nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("payload: invalid JSON: %w", err)
	}
	return &Builder{fields: m}, nil
}

// FromPairs parses a slice of "key=value" strings and merges them
// into the builder, overwriting any existing key.
func (b *Builder) FromPairs(pairs []string) error {
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("payload: invalid pair %q, expected key=value", p)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return fmt.Errorf("payload: empty key in pair %q", p)
		}
		b.fields[key] = strings.TrimSpace(parts[1])
	}
	return nil
}

// Build returns the assembled payload as a map.
func (b *Builder) Build() map[string]interface{} {
	out := make(map[string]interface{}, len(b.fields))
	for k, v := range b.fields {
		out[k] = v
	}
	return out
}

// JSON serialises the current payload to a JSON byte slice.
func (b *Builder) JSON() ([]byte, error) {
	return json.Marshal(b.fields)
}

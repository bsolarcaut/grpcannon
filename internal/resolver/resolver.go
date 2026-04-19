// Package resolver maps a gRPC method string to its fully-qualified form
// and validates that it conforms to the expected "/package.Service/Method" shape.
package resolver

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidMethod is returned when the method string cannot be parsed.
var ErrInvalidMethod = errors.New("resolver: invalid gRPC method format")

// Method holds the decomposed parts of a fully-qualified gRPC method.
type Method struct {
	PackageService string // e.g. "helloworld.Greeter"
	Name           string // e.g. "SayHello"
	Full           string // e.g. "/helloworld.Greeter/SayHello"
}

// Parse parses a gRPC method string into a Method.
// Accepted forms:
//   - /package.Service/Method
//   - package.Service/Method
func Parse(s string) (Method, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Method{}, fmt.Errorf("%w: empty string", ErrInvalidMethod)
	}

	norm := strings.TrimPrefix(s, "/")
	parts := strings.SplitN(norm, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Method{}, fmt.Errorf("%w: %q", ErrInvalidMethod, s)
	}

	return Method{
		PackageService: parts[0],
		Name:           parts[1],
		Full:           "/" + norm,
	}, nil
}

// MustParse is like Parse but panics on error.
func MustParse(s string) Method {
	m, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return m
}

// ServiceName returns only the unqualified service name (after the last dot).
func (m Method) ServiceName() string {
	idx := strings.LastIndex(m.PackageService, ".")
	if idx < 0 {
		return m.PackageService
	}
	return m.PackageService[idx+1:]
}

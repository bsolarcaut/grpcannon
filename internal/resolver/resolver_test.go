package resolver_test

import (
	"errors"
	"testing"

	"github.com/yourorg/grpcannon/internal/resolver"
)

func TestParse_ValidWithLeadingSlash(t *testing.T) {
	m, err := resolver.Parse("/helloworld.Greeter/SayHello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Full != "/helloworld.Greeter/SayHello" {
		t.Errorf("Full = %q", m.Full)
	}
	if m.PackageService != "helloworld.Greeter" {
		t.Errorf("PackageService = %q", m.PackageService)
	}
	if m.Name != "SayHello" {
		t.Errorf("Name = %q", m.Name)
	}
}

func TestParse_ValidWithoutLeadingSlash(t *testing.T) {
	m, err := resolver.Parse("helloworld.Greeter/SayHello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Full != "/helloworld.Greeter/SayHello" {
		t.Errorf("Full = %q", m.Full)
	}
}

func TestParse_EmptyStringReturnsError(t *testing.T) {
	_, err := resolver.Parse("")
	if !errors.Is(err, resolver.ErrInvalidMethod) {
		t.Fatalf("expected ErrInvalidMethod, got %v", err)
	}
}

func TestParse_MissingMethodPart(t *testing.T) {
	_, err := resolver.Parse("/helloworld.Greeter")
	if !errors.Is(err, resolver.ErrInvalidMethod) {
		t.Fatalf("expected ErrInvalidMethod, got %v", err)
	}
}

func TestParse_EmptyServicePart(t *testing.T) {
	_, err := resolver.Parse("/SayHello")
	if !errors.Is(err, resolver.ErrInvalidMethod) {
		t.Fatalf("expected ErrInvalidMethod, got %v", err)
	}
}

func TestServiceName_WithPackage(t *testing.T) {
	m := resolver.MustParse("/helloworld.Greeter/SayHello")
	if got := m.ServiceName(); got != "Greeter" {
		t.Errorf("ServiceName = %q, want Greeter", got)
	}
}

func TestServiceName_WithoutPackage(t *testing.T) {
	m := resolver.MustParse("/Greeter/SayHello")
	if got := m.ServiceName(); got != "Greeter" {
		t.Errorf("ServiceName = %q, want Greeter", got)
	}
}

func TestMustParse_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	resolver.MustParse("")
}

package metadata_test

import (
	"testing"

	"github.com/patrickward/grpcannon/internal/metadata"
)

func TestFromSlice_Empty(t *testing.T) {
	md, err := metadata.FromSlice(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(md) != 0 {
		t.Errorf("expected empty metadata, got %v", md)
	}
}

func TestFromSlice_ValidPairs(t *testing.T) {
	pairs := []string{"Authorization: Bearer token123", "x-request-id: abc"}
	md, err := metadata.FromSlice(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := md["authorization"]; len(got) != 1 || got[0] != "Bearer token123" {
		t.Errorf("authorization: got %v", got)
	}
	if got := md["x-request-id"]; len(got) != 1 || got[0] != "abc" {
		t.Errorf("x-request-id: got %v", got)
	}
}

func TestFromSlice_MissingColon(t *testing.T) {
	_, err := metadata.FromSlice([]string{"no-colon-here"})
	if err == nil {
		t.Fatal("expected error for missing colon, got nil")
	}
}

func TestFromSlice_EmptyKey(t *testing.T) {
	_, err := metadata.FromSlice([]string{": value"})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestFromSlice_DuplicateKeys(t *testing.T) {
	pairs := []string{"x-role: admin", "x-role: user"}
	md, err := metadata.FromSlice(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(md["x-role"]) != 2 {
		t.Errorf("expected 2 values for x-role, got %d", len(md["x-role"]))
	}
}

func TestMerge_CombinesTwoMaps(t *testing.T) {
	a, _ := metadata.FromSlice([]string{"x-a: 1"})
	b, _ := metadata.FromSlice([]string{"x-b: 2"})
	out := metadata.Merge(a, b)
	if len(out) != 2 {
		t.Errorf("expected 2 keys after merge, got %d", len(out))
	}
}

func TestMerge_BOverridesA(t *testing.T) {
	a, _ := metadata.FromSlice([]string{"x-key: original"})
	b, _ := metadata.FromSlice([]string{"x-key: overridden"})
	out := metadata.Merge(a, b)
	if got := out["x-key"]; len(got) != 1 || got[0] != "overridden" {
		t.Errorf("expected overridden value, got %v", got)
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	a, _ := metadata.FromSlice([]string{"x-a: 1"})
	b, _ := metadata.FromSlice([]string{"x-b: 2"})
	_ = metadata.Merge(a, b)
	if len(a) != 1 {
		t.Error("merge mutated input map a")
	}
}

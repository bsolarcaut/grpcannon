package payload_test

import (
	"encoding/json"
	"testing"

	"github.com/user/grpcannon/internal/payload"
)

func TestNew_EmptyBuild(t *testing.T) {
	b := payload.New()
	m := b.Build()
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}

func TestFromJSON_ValidObject(t *testing.T) {
	b, err := payload.FromJSON(`{"name":"alice","age":30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := b.Build()
	if m["name"] != "alice" {
		t.Errorf("expected name=alice, got %v", m["name"])
	}
}

func TestFromJSON_EmptyString(t *testing.T) {
	b, err := payload.FromJSON("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Build()) != 0 {
		t.Error("expected empty map for empty input")
	}
}

func TestFromJSON_InvalidJSON(t *testing.T) {
	_, err := payload.FromJSON(`{not valid}`)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestFromPairs_ValidPairs(t *testing.T) {
	b := payload.New()
	if err := b.FromPairs([]string{"foo=bar", "baz=qux"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := b.Build()
	if m["foo"] != "bar" || m["baz"] != "qux" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestFromPairs_MissingEquals(t *testing.T) {
	b := payload.New()
	if err := b.FromPairs([]string{"noequalssign"}); err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestFromPairs_EmptyKey(t *testing.T) {
	b := payload.New()
	if err := b.FromPairs([]string{"=value"}); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestFromPairs_OverwritesExisting(t *testing.T) {
	b := payload.New()
	_ = b.FromPairs([]string{"k=first"})
	_ = b.FromPairs([]string{"k=second"})
	if b.Build()["k"] != "second" {
		t.Error("expected key to be overwritten")
	}
}

func TestJSON_RoundTrip(t *testing.T) {
	b, _ := payload.FromJSON(`{"x":1}`)
	raw, err := b.JSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("round-trip unmarshal failed: %v", err)
	}
}

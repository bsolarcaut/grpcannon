package evict_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/evict"
)

func TestNew_CapClampedToOne(t *testing.T) {
	c := evict.New(0)
	c.Set("a", 1)
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}

func TestSet_Get_RoundTrip(t *testing.T) {
	c := evict.New(4)
	c.Set("k", "hello")
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v.(string) != "hello" {
		t.Fatalf("expected 'hello', got %v", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	c := evict.New(4)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for absent key")
	}
}

func TestSet_UpdateExisting(t *testing.T) {
	c := evict.New(4)
	c.Set("k", 1)
	c.Set("k", 2)
	v, _ := c.Get("k")
	if v.(int) != 2 {
		t.Fatalf("expected updated value 2, got %v", v)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len 1 after update, got %d", c.Len())
	}
}

func TestEviction_LRUDropped(t *testing.T) {
	c := evict.New(3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	// Access "a" so "b" becomes LRU.
	c.Get("a")
	// Adding "d" should evict "b".
	c.Set("d", 4)

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected 'b' to be evicted")
	}
	for _, k := range []string{"a", "c", "d"} {
		if _, ok := c.Get(k); !ok {
			t.Fatalf("expected key %q to be present", k)
		}
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	c := evict.New(4)
	c.Set("x", 99)
	c.Delete("x")
	if _, ok := c.Get("x"); ok {
		t.Fatal("expected key to be absent after delete")
	}
}

func TestDelete_NoopOnMissing(t *testing.T) {
	c := evict.New(4)
	c.Delete("ghost") // must not panic
}

func TestConcurrentAccess(t *testing.T) {
	c := evict.New(64)
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i%32)
			c.Set(key, i)
			c.Get(key)
		}(i)
	}
	wg.Wait()
	if c.Len() > 64 {
		t.Fatalf("cache exceeded capacity: len=%d", c.Len())
	}
}

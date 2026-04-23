// Package evict implements a fixed-capacity LRU cache used to deduplicate
// in-flight request keys and evict stale entries under memory pressure.
package evict

import (
	"container/list"
	"sync"
)

// entry holds a key/value pair stored in the cache.
type entry struct {
	key   string
	value any
}

// Cache is a thread-safe LRU cache with a fixed capacity.
type Cache struct {
	mu       sync.Mutex
	cap      int
	items    map[string]*list.Element
	order    *list.List
}

// New returns a Cache that holds at most cap entries.
// If cap is less than 1 it is clamped to 1.
func New(cap int) *Cache {
	if cap < 1 {
		cap = 1
	}
	return &Cache{
		cap:   cap,
		items: make(map[string]*list.Element, cap),
		order: list.New(),
	}
}

// Set inserts or updates key with value, evicting the least-recently-used
// entry when the cache is at capacity.
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.order.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	if c.order.Len() >= c.cap {
		c.evictLRU()
	}

	el := c.order.PushFront(&entry{key: key, value: value})
	c.items[key] = el
}

// Get returns the value associated with key and whether it was found.
// A successful lookup promotes the entry to most-recently-used.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry).value, true
}

// Delete removes key from the cache. It is a no-op if key is absent.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.removeElement(el)
	}
}

// Len returns the current number of entries in the cache.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

func (c *Cache) evictLRU() {
	el := c.order.Back()
	if el != nil {
		c.removeElement(el)
	}
}

func (c *Cache) removeElement(el *list.Element) {
	c.order.Remove(el)
	delete(c.items, el.Value.(*entry).key)
}

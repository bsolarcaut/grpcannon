// Package evict provides a thread-safe, fixed-capacity LRU (Least Recently
// Used) cache suitable for deduplicating short-lived request keys during a
// load-test run.
//
// Usage:
//
//	c := evict.New(512)
//	c.Set("req-id-1", struct{}{})
//	if _, ok := c.Get("req-id-1"); ok {
//		// duplicate – skip
//	}
//
// The cache is safe for concurrent use by multiple goroutines.
package evict

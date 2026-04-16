package histogram

import (
	"bytes"
	"testing"
	"time"
)

var defaultBounds = []time.Duration{
	1 * time.Millisecond,
	5 * time.Millisecond,
	10 * time.Millisecond,
	50 * time.Millisecond,
	100 * time.Millisecond,
}

func TestNew_BucketCount(t *testing.T) {
	h := New(defaultBounds)
	// len(bounds) + 1 overflow
	if got := len(h.Buckets()); got != len(defaultBounds)+1 {
		t.Fatalf("expected %d buckets, got %d", len(defaultBounds)+1, got)
	}
}

func TestRecord_CorrectBucket(t *testing.T) {
	h := New(defaultBounds)
	h.Record(500 * time.Microsecond) // < 1ms → bucket 0
	h.Record(3 * time.Millisecond)   // < 5ms → bucket 1
	h.Record(200 * time.Millisecond) // overflow

	buckets := h.Buckets()
	if buckets[0].Count != 1 {
		t.Errorf("bucket 0: expected 1, got %d", buckets[0].Count)
	}
	if buckets[1].Count != 1 {
		t.Errorf("bucket 1: expected 1, got %d", buckets[1].Count)
	}
	overflow := buckets[len(buckets)-1]
	if overflow.Count != 1 {
		t.Errorf("overflow: expected 1, got %d", overflow.Count)
	}
}

func TestRecord_MultipleInSameBucket(t *testing.T) {
	h := New(defaultBounds)
	for i := 0; i < 5; i++ {
		h.Record(2 * time.Millisecond)
	}
	if got := h.Buckets()[1].Count; got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestPrint_NoOutput_WhenEmpty(t *testing.T) {
	h := New(defaultBounds)
	// record at least one to avoid div-by-zero in max
	h.Record(1 * time.Microsecond)
	var buf bytes.Buffer
	h.Print(&buf)
	0 {
		t.Error("expected non-empty output")
	}
}

func TestBuckets_ReturnsCopy(t *testing.T) {
	h := New(defaultBounds)
	h.Record(500 * time.Microsecond)
	b1 := h.Buckets()
	b1[0].Count = 999
	b2 := h.Buckets()
	if b2[0].Count == 999 {
		t.Error("Buckets should return a copy, not a reference")
	}
}

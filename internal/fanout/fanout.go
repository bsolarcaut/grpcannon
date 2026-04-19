// Package fanout distributes a stream of results to multiple consumers.
package fanout

import "sync"

// Sink is a channel that receives result values.
type Sink[T any] <-chan T

// Fanout fans a single input channel out to N subscriber channels.
type Fanout[T any] struct {
	mu   sync.Mutex
	subs []chan T
	buf  int
}

// New creates a Fanout with the given per-subscriber buffer size.
func New[T any](buf int) *Fanout[T] {
	if buf < 1 {
		buf = 1
	}
	return &Fanout[T]{buf: buf}
}

// Subscribe returns a new Sink that will receive every value sent to Run.
func (f *Fanout[T]) Subscribe() Sink[T] {
	ch := make(chan T, f.buf)
	f.mu.Lock()
	f.subs = append(f.subs, ch)
	f.mu.Unlock()
	return ch
}

// Run reads from src and broadcasts each value to all subscribers.
// It closes all subscriber channels when src is closed.
func (f *Fanout[T]) Run(src <-chan T) {
	for v := range src {
		f.mu.Lock()
		for _, ch := range f.subs {
			ch <- v
		}
		f.mu.Unlock()
	}
	f.mu.Lock()
	for _, ch := range f.subs {
		close(ch)
	}
	f.mu.Unlock()
}

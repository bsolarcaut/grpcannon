// Package debounce provides a debouncer that delays execution of a function
// until after a specified quiet period has elapsed since the last invocation.
// Useful for coalescing bursts of events (e.g. config reloads, metric flushes)
// into a single action.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays calls to fn until no new calls have arrived within wait.
type Debouncer struct {
	wait  time.Duration
	fn    func()
	mu    sync.Mutex
	timer *time.Timer
}

// New returns a Debouncer that will call fn after wait has elapsed since the
// last call to Trigger. A zero or negative wait is clamped to 1ms.
func New(wait time.Duration, fn func()) *Debouncer {
	if wait <= 0 {
		wait = time.Millisecond
	}
	return &Debouncer{wait: wait, fn: fn}
}

// Trigger schedules fn to be called after the debounce window. If Trigger is
// called again before the window expires the timer resets.
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Reset(d.wait)
		return
	}
	d.timer = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		d.timer = nil
		d.mu.Unlock()
		d.fn()
	})
}

// Flush cancels any pending timer and invokes fn immediately. If no call is
// pending Flush is a no-op.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	if d.timer == nil {
		d.mu.Unlock()
		return
	}
	d.timer.Stop()
	d.timer = nil
	d.mu.Unlock()
	d.fn()
}

// Stop cancels any pending timer without invoking fn.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}

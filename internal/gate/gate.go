// Package gate provides a simple boolean gate that can be used to
// pause and resume request flow during a load test.
package gate

import "sync"

// Gate is a pausable flow-control primitive. When closed, callers
// block on Wait until the gate is opened again.
type Gate struct {
	mu     sync.Mutex
	cond   *sync.Cond
	open   bool
	closed bool // permanently closed (stopped)
}

// New returns an open Gate.
func New() *Gate {
	g := &Gate{open: true}
	g.cond = sync.NewCond(&g.mu)
	return g
}

// Pause closes the gate; subsequent calls to Wait will block.
func (g *Gate) Pause() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.open = false
}

// Resume opens the gate and unblocks all waiting callers.
func (g *Gate) Resume() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.open = true
	g.cond.Broadcast()
}

// Stop permanently opens the gate so no caller blocks forever.
func (g *Gate) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.closed = true
	g.open = true
	g.cond.Broadcast()
}

// Wait blocks until the gate is open or stopped.
func (g *Gate) Wait() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for !g.open && !g.closed {
		g.cond.Wait()
	}
}

// IsOpen reports whether the gate is currently open.
func (g *Gate) IsOpen() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.open
}

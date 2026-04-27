// Package shedder implements adaptive load shedding based on a configurable
// error-rate threshold and a sliding window of recent outcomes.
// When the error rate exceeds the threshold the shedder rejects new requests
// until the window recovers.
package shedder

import (
	"errors"
	"sync"
	"time"
)

// ErrShed is returned when the shedder rejects a request.
var ErrShed = errors.New("shedder: request shed")

// Config holds tunable parameters for the Shedder.
type Config struct {
	// Threshold is the error-rate fraction [0,1] above which requests are shed.
	Threshold float64
	// WindowSize is the number of outcomes tracked in the sliding window.
	WindowSize int
	// CoolDown is the minimum time to wait before re-evaluating after shedding.
	CoolDown time.Duration
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Threshold:  0.5,
		WindowSize: 100,
		CoolDown:   5 * time.Second,
	}
}

// Shedder tracks recent outcomes and sheds load when the error rate is too high.
type Shedder struct {
	cfg      Config
	mu       sync.Mutex
	bucket   []bool // true == error
	head     int
	filled   int
	shedding bool
	shedAt   time.Time
}

// New creates a Shedder using the provided Config.
func New(cfg Config) *Shedder {
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = Default().WindowSize
	}
	if cfg.Threshold <= 0 || cfg.Threshold > 1 {
		cfg.Threshold = Default().Threshold
	}
	if cfg.CoolDown <= 0 {
		cfg.CoolDown = Default().CoolDown
	}
	return &Shedder{
		cfg:    cfg,
		bucket: make([]bool, cfg.WindowSize),
	}
}

// Allow returns nil when the request may proceed, or ErrShed when it should be
// dropped. It must be paired with a call to Record.
func (s *Shedder) Allow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.shedding {
		if time.Since(s.shedAt) < s.cfg.CoolDown {
			return ErrShed
		}
		s.shedding = false
	}
	return nil
}

// Record registers an outcome. isErr should be true when the call failed.
func (s *Shedder) Record(isErr bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bucket[s.head] = isErr
	s.head = (s.head + 1) % s.cfg.WindowSize
	if s.filled < s.cfg.WindowSize {
		s.filled++
	}
	if s.filled == 0 {
		return
	}
	var errs int
	for i := 0; i < s.filled; i++ {
		if s.bucket[i] {
			errs++
		}
	}
	rate := float64(errs) / float64(s.filled)
	if rate >= s.cfg.Threshold {
		s.shedding = true
		s.shedAt = time.Now()
	}
}

// Rate returns the current error rate in the window [0,1].
func (s *Shedder) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.filled == 0 {
		return 0
	}
	var errs int
	for i := 0; i < s.filled; i++ {
		if s.bucket[i] {
			errs++
		}
	}
	return float64(errs) / float64(s.filled)
}

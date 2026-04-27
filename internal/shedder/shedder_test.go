package shedder_test

import (
	"testing"
	"time"

	"github.com/yourorg/grpcannon/internal/shedder"
)

func TestDefault_Fields(t *testing.T) {
	cfg := shedder.Default()
	if cfg.Threshold != 0.5 {
		t.Fatalf("expected threshold 0.5, got %v", cfg.Threshold)
	}
	if cfg.WindowSize != 100 {
		t.Fatalf("expected window 100, got %d", cfg.WindowSize)
	}
	if cfg.CoolDown != 5*time.Second {
		t.Fatalf("expected cooldown 5s, got %v", cfg.CoolDown)
	}
}

func TestAllow_InitiallyPermits(t *testing.T) {
	s := shedder.New(shedder.Default())
	if err := s.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecord_BelowThreshold_Allows(t *testing.T) {
	cfg := shedder.Default()
	cfg.WindowSize = 10
	cfg.Threshold = 0.6
	s := shedder.New(cfg)
	for i := 0; i < 4; i++ {
		s.Record(true)
	}
	for i := 0; i < 6; i++ {
		s.Record(false)
	}
	if err := s.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecord_ExceedsThreshold_Sheds(t *testing.T) {
	cfg := shedder.Default()
	cfg.WindowSize = 10
	cfg.Threshold = 0.5
	cfg.CoolDown = 10 * time.Second // long cooldown so it stays shed
	s := shedder.New(cfg)
	for i := 0; i < 6; i++ {
		s.Record(true)
	}
	for i := 0; i < 4; i++ {
		s.Record(false)
	}
	if err := s.Allow(); err != shedder.ErrShed {
		t.Fatalf("expected ErrShed, got %v", err)
	}
}

func TestRate_Empty(t *testing.T) {
	s := shedder.New(shedder.Default())
	if r := s.Rate(); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRate_AllErrors(t *testing.T) {
	cfg := shedder.Default()
	cfg.WindowSize = 4
	s := shedder.New(cfg)
	for i := 0; i < 4; i++ {
		s.Record(true)
	}
	if r := s.Rate(); r != 1.0 {
		t.Fatalf("expected 1.0, got %v", r)
	}
}

func TestAllow_RecoverAfterCooldown(t *testing.T) {
	cfg := shedder.Default()
	cfg.WindowSize = 4
	cfg.Threshold = 0.5
	cfg.CoolDown = 1 * time.Millisecond
	s := shedder.New(cfg)
	for i := 0; i < 4; i++ {
		s.Record(true)
	}
	if err := s.Allow(); err != shedder.ErrShed {
		t.Fatalf("expected ErrShed immediately after threshold, got %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	if err := s.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
}

func TestNew_ClampsInvalidWindowSize(t *testing.T) {
	cfg := shedder.Config{WindowSize: 0, Threshold: 0.5, CoolDown: time.Second}
	s := shedder.New(cfg)
	if s == nil {
		t.Fatal("expected non-nil shedder")
	}
}

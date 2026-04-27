package quorum_test

import (
	"errors"
	"testing"

	"github.com/example/grpcannon/internal/quorum"
)

func TestNew_DefaultThresholdOnZero(t *testing.T) {
	q := quorum.New(0)
	q.Vote(true)
	q.Vote(false) // 50 % healthy — meets default 0.5
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNew_ClampsAboveOne(t *testing.T) {
	q := quorum.New(2.0)
	q.Vote(true)
	if err := q.Check(); err != nil {
		t.Fatalf("threshold clamped to 1.0 should pass when all healthy: %v", err)
	}
}

func TestCheck_NoVotesReturnsNil(t *testing.T) {
	q := quorum.New(0.8)
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil with no votes, got %v", err)
	}
}

func TestCheck_AllHealthy(t *testing.T) {
	q := quorum.New(0.75)
	for i := 0; i < 4; i++ {
		q.Vote(true)
	}
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_BelowThreshold(t *testing.T) {
	q := quorum.New(0.75)
	q.Vote(true)
	q.Vote(false)
	q.Vote(false)
	q.Vote(false) // 25 % healthy
	err := q.Check()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, quorum.ErrQuorumNotMet) {
		t.Fatalf("expected ErrQuorumNotMet, got %v", err)
	}
}

func TestCheck_ExactlyAtThreshold(t *testing.T) {
	q := quorum.New(0.5)
	q.Vote(true)
	q.Vote(false) // exactly 50 %
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil at exact threshold, got %v", err)
	}
}

func TestReset_ClearsVotes(t *testing.T) {
	q := quorum.New(0.9)
	q.Vote(false)
	q.Vote(false)
	q.Reset()
	if err := q.Check(); err != nil {
		t.Fatalf("after reset expected nil, got %v", err)
	}
}

func TestFraction_NoVotes(t *testing.T) {
	q := quorum.New(0.5)
	if f := q.Fraction(); f != 0 {
		t.Fatalf("expected 0, got %f", f)
	}
}

func TestFraction_PartialHealthy(t *testing.T) {
	q := quorum.New(0.5)
	q.Vote(true)
	q.Vote(true)
	q.Vote(false)
	q.Vote(false)
	if got := q.Fraction(); got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

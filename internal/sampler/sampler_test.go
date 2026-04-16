package sampler

import (
	"testing"
)

func TestNew_ClampsBelowZero(t *testing.T) {
	s := New(-0.5, 42)
	if s.Rate() != 0 {
		t.Fatalf("expected rate 0, got %f", s.Rate())
	}
}

func TestNew_ClampsAboveOne(t *testing.T) {
	s := New(1.5, 42)
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %f", s.Rate())
	}
}

func TestSample_ZeroRateNeverSamples(t *testing.T) {
	s := New(0, 42)
	for i := 0; i < 1000; i++ {
		if s.Sample() {
			t.Fatal("expected no samples with rate=0")
		}
	}
	total, sampled := s.Stats()
	if total != 1000 || sampled != 0 {
		t.Fatalf("unexpected stats: total=%d sampled=%d", total, sampled)
	}
}

func TestSample_FullRateSamplesAll(t *testing.T) {
	s := New(1.0, 42)
	for i := 0; i < 500; i++ {
		if !s.Sample() {
			t.Fatal("expected all samples with rate=1")
		}
	}
	total, sampled := s.Stats()
	if total != 500 || sampled != 500 {
		t.Fatalf("unexpected stats: total=%d sampled=%d", total, sampled)
	}
}

func TestSample_PartialRate_ApproximatelyCorrect(t *testing.T) {
	const n = 100_000
	s := New(0.1, 99)
	for i := 0; i < n; i++ {
		s.Sample()
	}
	total, sampled := s.Stats()
	if total != n {
		t.Fatalf("total mismatch: got %d", total)
	}
	ratio := float64(sampled) / float64(total)
	if ratio < 0.08 || ratio > 0.12 {
		t.Fatalf("sample ratio %f out of expected range [0.08, 0.12]", ratio)
	}
}

func TestStats_InitiallyZero(t *testing.T) {
	s := New(0.5, 1)
	total, sampled := s.Stats()
	if total != 0 || sampled != 0 {
		t.Fatalf("expected zero stats, got total=%d sampled=%d", total, sampled)
	}
}

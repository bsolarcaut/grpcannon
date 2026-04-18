package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/grpcannon/internal/profile"
)

func TestStartCPU_EmptyPathIsNoop(t *testing.T) {
	p, err := profile.StartCPU("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}
}

func TestStartCPU_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cpu.prof")

	p, err := profile.StartCPU(path)
	if err != nil {
		t.Fatalf("start error: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected cpu profile file: %v", err)
	}
}

func TestWithMemPath_WritesFile(t *testing.T) {
	dir := t.TempDir()
	memPath := filepath.Join(dir, "mem.prof")

	p, err := profile.StartCPU("")
	if err != nil {
		t.Fatalf("start error: %v", err)
	}
	p.WithMemPath(memPath)
	if err := p.Stop(); err != nil {
		t.Fatalf("stop error: %v", err)
	}
	if _, err := os.Stat(memPath); err != nil {
		t.Fatalf("expected mem profile file: %v", err)
	}
}

func TestStop_Idempotent(t *testing.T) {
	p, err := profile.StartCPU("")
	if err != nil {
		t.Fatalf("start error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := p.Stop(); err != nil {
			t.Fatalf("stop %d error: %v", i, err)
		}
	}
}

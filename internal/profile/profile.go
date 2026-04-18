// Package profile provides CPU and memory profiling helpers for grpcannon.
package profile

import (
	"fmt"
	"os"
	"runtime/pprof"
)

// Profile holds open file handles for active profiles.
type Profile struct {
	cpuFile *os.File
	memPath string
}

// StartCPU begins CPU profiling, writing to path.
func StartCPU(path string) (*Profile, error) {
	if path == "" {
		return &Profile{}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("profile: create cpu file: %w", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("profile: start cpu profile: %w", err)
	}
	return &Profile{cpuFile: f}, nil
}

// WithMemPath sets the path for a memory profile written on Stop.
func (p *Profile) WithMemPath(path string) *Profile {
	p.memPath = path
	return p
}

// Stop finalises all active profiles.
func (p *Profile) Stop() error {
	if p.cpuFile != nil {
		pprof.StopCPUProfile()
		if err := p.cpuFile.Close(); err != nil {
			return fmt.Errorf("profile: close cpu file: %w", err)
		}
		p.cpuFile = nil
	}
	if p.memPath != "" {
		f, err := os.Create(p.memPath)
		if err != nil {
			return fmt.Errorf("profile: create mem file: %w", err)
		}
		defer f.Close()
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("profile: write heap profile: %w", err)
		}
	}
	return nil
}

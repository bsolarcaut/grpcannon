package tui_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nickcorin/grpcannon/internal/tui"
)

func TestBanner_ContainsTitle(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, false)
	p.Banner()
	if !strings.Contains(buf.String(), "grpcannon") {
		t.Fatalf("expected banner to contain 'grpcannon', got: %q", buf.String())
	}
}

func TestProgress_ShowsPercentage(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, false)
	p.Progress(50, 100, 2*time.Second)
	out := buf.String()
	if !strings.Contains(out, "50%") {
		t.Fatalf("expected 50%% in output, got: %q", out)
	}
}

func TestProgress_ZeroTotal(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, false)
	p.Progress(0, 0, 0)
	if !strings.Contains(buf.String(), "0%") {
		t.Fatal("expected 0% for zero total")
	}
}

func TestVerbose_PrintsWhenEnabled(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, true)
	p.Verbose("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("expected verbose output, got: %q", buf.String())
	}
}

func TestVerbose_SilentWhenDisabled(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, false)
	p.Verbose("should not appear")
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %q", buf.String())
	}
}

func TestError_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := tui.New(&buf, false)
	p.Error(errors.New("something failed"))
	if !strings.Contains(buf.String(), "something failed") {
		t.Fatalf("expected error message, got: %q", buf.String())
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic.
	p := tui.New(nil, false)
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
}

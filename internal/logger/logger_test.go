package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_NilWriterDefaultsToStderr(t *testing.T) {
	l := New(nil, LevelInfo)
	if l.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestLog_BelowLevelSuppressed(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelInfo)
	l.Debug("should not appear")
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestLog_AtLevelWritten(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelInfo)
	l.Info("hello %s", "world")
	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected [INFO] in output, got %q", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected message in output, got %q", out)
	}
}

func TestLog_AboveLevelWritten(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelWarn)
	l.Error("something broke")
	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got %q", out)
	}
}

func TestLog_WarnSuppressedBelowError(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelError)
	l.Warn("this is a warning")
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestDefault_ReturnsNonNil(t *testing.T) {
	if Default() == nil {
		t.Fatal("expected non-nil default logger")
	}
}

func TestLog_FormatsArgs(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, LevelDebug)
	l.Debug("count=%d", 42)
	if !strings.Contains(buf.String(), "count=42") {
		t.Errorf("expected formatted message, got %q", buf.String())
	}
}

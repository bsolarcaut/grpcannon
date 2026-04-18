// Package logger provides a lightweight structured logger for grpcannon.
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a log message.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

// Logger writes structured log lines to an io.Writer.
type Logger struct {
	w     io.Writer
	level Level
}

// New returns a Logger that writes to w at the given minimum level.
func New(w io.Writer, level Level) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w, level: level}
}

// Default returns a Logger writing to stderr at Info level.
func Default() *Logger { return New(os.Stderr, LevelInfo) }

func (l *Logger) log(level Level, msg string, args ...any) {
	if level < l.level {
		return
	}
	line := fmt.Sprintf("%s [%s] %s", time.Now().UTC().Format(time.RFC3339), levelNames[level], fmt.Sprintf(msg, args...))
	fmt.Fprintln(l.w, line)
}

// Debug logs at DEBUG level.
func (l *Logger) Debug(msg string, args ...any) { l.log(LevelDebug, msg, args...) }

// Info logs at INFO level.
func (l *Logger) Info(msg string, args ...any) { l.log(LevelInfo, msg, args...) }

// Warn logs at WARN level.
func (l *Logger) Warn(msg string, args ...any) { l.log(LevelWarn, msg, args...) }

// Error logs at ERROR level.
func (l *Logger) Error(msg string, args ...any) { l.log(LevelError, msg, args...) }

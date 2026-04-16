// Package tui provides terminal output helpers for grpcannon.
package tui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Printer writes progress and status lines to a writer.
type Printer struct {
	w       io.Writer
	verbose bool
}

// New returns a Printer writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer, verbose bool) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w, verbose: verbose}
}

// Banner prints the grpcannon banner.
func (p *Printer) Banner() {
	fmt.Fprintln(p.w, strings.Repeat("=", 40))
	fmt.Fprintln(p.w, "  grpcannon — gRPC load tester")
	fmt.Fprintln(p.w, strings.Repeat("=", 40))
}

// Progress prints a single progress line with requests sent and elapsed time.
func (p *Printer) Progress(sent, total int, elapsed time.Duration) {
	pct := 0
	if total > 0 {
		pct = sent * 100 / total
	}
	fmt.Fprintf(p.w, "\r  Sent: %d / %d (%d%%)  Elapsed: %s",
		sent, total, pct, elapsed.Round(time.Millisecond))
}

// Done prints a newline after progress output.
func (p *Printer) Done() {
	fmt.Fprintln(p.w)
}

// Verbose prints msg only when verbose mode is enabled.
func (p *Printer) Verbose(format string, args ...any) {
	if p.verbose {
		fmt.Fprintf(p.w, "[verbose] "+format+"\n", args...)
	}
}

// Error prints an error message.
func (p *Printer) Error(err error) {
	fmt.Fprintf(p.w, "[error] %v\n", err)
}

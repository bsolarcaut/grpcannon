package observe_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nickpoorman/grpcannon/internal/observe"
	"github.com/nickpoorman/grpcannon/internal/snapshot"
)

func TestObserver_EmitsLines(t *testing.T) {
	c := snapshot.NewCollector()
	c.RecordSuccess(10 * time.Millisecond)
	c.RecordSuccess(20 * time.Millisecond)
	c.RecordError(5 * time.Millisecond)

	var buf bytes.Buffer
	obs := observe.New(c,
		observe.WithInterval(20*time.Millisecond),
		observe.WithWriter(&buf),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	obs.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "[observe]") {
		t.Fatalf("expected observe prefix, got: %q", out)
	}
	if !strings.Contains(out, "total=3") {
		t.Fatalf("expected total=3 in output, got: %q", out)
	}
	if !strings.Contains(out, "err=1") {
		t.Fatalf("expected err=1 in output, got: %q", out)
	}
}

func TestObserver_StopsOnCancel(t *testing.T) {
	c := snapshot.NewCollector()
	var buf bytes.Buffer
	obs := observe.New(c,
		observe.WithInterval(10*time.Millisecond),
		observe.WithWriter(&buf),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		obs.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("observer did not stop after context cancellation")
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	c := snapshot.NewCollector()
	obs := observe.New(c)
	if obs == nil {
		t.Fatal("expected non-nil observer")
	}
}

func TestWithInterval_ZeroIsIgnored(t *testing.T) {
	c := snapshot.NewCollector()
	var buf bytes.Buffer
	obs := observe.New(c, observe.WithInterval(0), observe.WithWriter(&buf))
	if obs == nil {
		t.Fatal("expected non-nil observer")
	}
}

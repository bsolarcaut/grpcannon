package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/retry"
)

var errTransient = errors.New("transient error")

func TestDefaultPolicy_Fields(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Fatalf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Delay != 100*time.Millisecond {
		t.Fatalf("unexpected delay: %v", p.Delay)
	}
}

func TestIsRetryable_NilError(t *testing.T) {
	if retry.IsRetryable(nil) {
		t.Fatal("nil error should not be retryable")
	}
}

func TestIsRetryable_ContextCanceled(t *testing.T) {
	if retry.IsRetryable(context.Canceled) {
		t.Fatal("context.Canceled should not be retryable")
	}
}

func TestIsRetryable_TransientError(t *testing.T) {
	if !retry.IsRetryable(errTransient) {
		t.Fatal("transient error should be retryable")
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	calls := 0
	err := p.Do(context.Background(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("expected 1 call and no error, got calls=%d err=%v", calls, err)
	}
}

func TestDo_RetriesOnTransientError(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0, RetryOn: retry.IsRetryable}
	calls := 0
	err := p.Do(context.Background(), func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, Delay: 0, RetryOn: retry.IsRetryable}
	calls := 0
	err := p.Do(context.Background(), func(_ context.Context) error {
		calls++
		return errTransient
	})
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected transient error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDo_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := retry.Policy{MaxAttempts: 5, Delay: 0, RetryOn: retry.IsRetryable}
	err := p.Do(ctx, func(_ context.Context) error { return errTransient })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

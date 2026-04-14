package timeout_test

import (
	"errors"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/timeout"
)

func fastFn(results []drift.Result, err error) timeout.RunFunc {
	return func() ([]drift.Result, error) {
		return results, err
	}
}

func slowFn(d time.Duration) timeout.RunFunc {
	return func() ([]drift.Result, error) {
		time.Sleep(d)
		return nil, nil
	}
}

func TestApply_NoDeadline_CallsFnDirectly(t *testing.T) {
	want := []drift.Result{{Service: "svc", Drifted: false}}
	opts := timeout.Options{Deadline: 0}
	got, err := timeout.Apply(opts, fastFn(want, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Service != "svc" {
		t.Fatalf("unexpected results: %v", got)
	}
}

func TestApply_FastFn_ReturnsResults(t *testing.T) {
	want := []drift.Result{{Service: "api", Drifted: true}}
	opts := timeout.DefaultOptions()
	got, err := timeout.Apply(opts, fastFn(want, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Service != "api" {
		t.Fatalf("unexpected results: %v", got)
	}
}

func TestApply_FastFn_PropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	opts := timeout.DefaultOptions()
	_, err := timeout.Apply(opts, fastFn(nil, sentinel))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestApply_SlowFn_ReturnsErrTimeout(t *testing.T) {
	opts := timeout.Options{
		Deadline:    20 * time.Millisecond,
		GracePeriod: 5 * time.Millisecond,
	}
	_, err := timeout.Apply(opts, slowFn(200*time.Millisecond))
	if !errors.Is(err, timeout.ErrTimeout) {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := timeout.DefaultOptions()
	if opts.Deadline != 30*time.Second {
		t.Errorf("expected 30s deadline, got %v", opts.Deadline)
	}
	if opts.GracePeriod != 5*time.Second {
		t.Errorf("expected 5s grace, got %v", opts.GracePeriod)
	}
}

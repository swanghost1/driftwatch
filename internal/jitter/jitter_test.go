package jitter_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/jitter"
)

func TestApply_ZeroFactor_ReturnsBase(t *testing.T) {
	base := 10 * time.Second
	opts := jitter.Options{Factor: 0}
	got := jitter.Apply(base, opts)
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_PositiveFactor_ResultAtLeastBase(t *testing.T) {
	base := 10 * time.Second
	opts := jitter.Options{Factor: 0.2, Seed: 42}
	got := jitter.Apply(base, opts)
	if got < base {
		t.Fatalf("expected result >= base %v, got %v", base, got)
	}
}

func TestApply_PositiveFactor_ResultWithinBounds(t *testing.T) {
	base := 10 * time.Second
	opts := jitter.Options{Factor: 0.2, Seed: 42}
	max := base + time.Duration(float64(base)*0.2)
	for i := 0; i < 50; i++ {
		opts.Seed = int64(i + 1)
		got := jitter.Apply(base, opts)
		if got > max {
			t.Fatalf("result %v exceeds max %v", got, max)
		}
	}
}

func TestApply_DefaultOptions_NonZeroFactor(t *testing.T) {
	opts := jitter.DefaultOptions()
	if opts.Factor == 0 {
		t.Fatal("expected non-zero default factor")
	}
}

func TestApplyFull_ZeroFactor_ReturnsBase(t *testing.T) {
	base := 5 * time.Second
	opts := jitter.Options{Factor: 0}
	got := jitter.ApplyFull(base, opts)
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApplyFull_WithFactor_ResultWithinSymmetricBounds(t *testing.T) {
	base := 10 * time.Second
	factor := 0.3
	opts := jitter.Options{Factor: factor}
	margin := time.Duration(float64(base) * factor)
	for i := 1; i <= 50; i++ {
		opts.Seed = int64(i)
		got := jitter.ApplyFull(base, opts)
		if got < base-margin || got > base+margin {
			t.Fatalf("result %v out of [%v, %v]", got, base-margin, base+margin)
		}
	}
}

func TestApply_Deterministic_SameSeed(t *testing.T) {
	base := 8 * time.Second
	opts := jitter.Options{Factor: 0.25, Seed: 99}
	a := jitter.Apply(base, opts)
	b := jitter.Apply(base, opts)
	if a != b {
		t.Fatalf("expected same result for same seed: %v vs %v", a, b)
	}
}

package sample_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/sample"
)

func makeResults(driftedCount, cleanCount int) []drift.Result {
	out := make([]drift.Result, 0, driftedCount+cleanCount)
	for i := 0; i < driftedCount; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: true, Field: "image"})
	}
	for i := 0; i < cleanCount; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: false})
	}
	return out
}

func TestShouldRecord_ZeroRate_ReturnsFalse(t *testing.T) {
	opts := sample.Options{Rate: 0.0, AlwaysSampleDrifted: false}
	results := makeResults(1, 1)
	if sample.ShouldRecord(results, opts) {
		t.Fatal("expected false for zero rate")
	}
}

func TestShouldRecord_FullRate_ReturnsTrue(t *testing.T) {
	opts := sample.Options{Rate: 1.0, AlwaysSampleDrifted: false}
	results := makeResults(0, 3)
	if !sample.ShouldRecord(results, opts) {
		t.Fatal("expected true for rate 1.0")
	}
}

func TestShouldRecord_AlwaysSampleDrifted_ReturnsTrueWhenDrifted(t *testing.T) {
	opts := sample.Options{Rate: 0.0, AlwaysSampleDrifted: true}
	results := makeResults(1, 0)
	if !sample.ShouldRecord(results, opts) {
		t.Fatal("expected true: drifted result with AlwaysSampleDrifted")
	}
}

func TestShouldRecord_AlwaysSampleDrifted_NoDrift_ReturnsFalse(t *testing.T) {
	opts := sample.Options{Rate: 0.0, AlwaysSampleDrifted: true}
	results := makeResults(0, 3)
	if sample.ShouldRecord(results, opts) {
		t.Fatal("expected false: no drift and rate is 0")
	}
}

func TestApply_NoMax_ReturnsAll(t *testing.T) {
	results := makeResults(2, 3)
	out := sample.Apply(results, 0)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_MaxLargerThanInput_ReturnsAll(t *testing.T) {
	results := makeResults(1, 2)
	out := sample.Apply(results, 100)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_TruncatesToN(t *testing.T) {
	results := makeResults(2, 5)
	out := sample.Apply(results, 3)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestApply_PrioritisesDrifted(t *testing.T) {
	results := makeResults(2, 5)
	out := sample.Apply(results, 2)
	for _, r := range out {
		if !r.Drifted {
			t.Fatal("expected only drifted results when n equals drifted count")
		}
	}
}

func TestDefaultOptions_SensibleDefaults(t *testing.T) {
	opts := sample.DefaultOptions()
	if opts.Rate <= 0 || opts.Rate > 1 {
		t.Fatalf("unexpected default rate: %f", opts.Rate)
	}
	if !opts.AlwaysSampleDrifted {
		t.Fatal("expected AlwaysSampleDrifted to be true by default")
	}
}

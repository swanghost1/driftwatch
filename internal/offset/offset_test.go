package offset_test

import (
	"bytes"
	"testing"
	"time"

	"driftwatch/internal/drift"
	"driftwatch/internal/offset"
)

func makeResults(drifted bool, t time.Time) drift.Result {
	return drift.Result{
		Service:    "svc",
		Field:      "image",
		Drifted:    drifted,
		DetectedAt: t,
	}
}

func TestApply_ZeroShift_ReturnsOriginal(t *testing.T) {
	now := time.Now()
	results := []drift.Result{makeResults(true, now)}
	out := offset.Apply(results, offset.DefaultOptions())
	if &out[0] == &results[0] {
		t.Fatal("expected same slice to be returned for zero shift")
	}
	if !out[0].DetectedAt.Equal(now) {
		t.Errorf("expected time unchanged, got %v", out[0].DetectedAt)
	}
}

func TestApply_PositiveShift_AdvancesTime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	results := []drift.Result{makeResults(true, now)}
	opts := offset.Options{Shift: 2 * time.Hour}
	out := offset.Apply(results, opts)
	want := now.Add(2 * time.Hour)
	if !out[0].DetectedAt.Equal(want) {
		t.Errorf("expected %v, got %v", want, out[0].DetectedAt)
	}
}

func TestApply_NegativeShift_RewindsTime(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	results := []drift.Result{makeResults(true, now)}
	opts := offset.Options{Shift: -30 * time.Minute}
	out := offset.Apply(results, opts)
	want := now.Add(-30 * time.Minute)
	if !out[0].DetectedAt.Equal(want) {
		t.Errorf("expected %v, got %v", want, out[0].DetectedAt)
	}
}

func TestApply_SkipClean_OnlyAdjustsDrifted(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	results := []drift.Result{
		makeResults(true, now),
		makeResults(false, now),
	}
	opts := offset.Options{Shift: time.Hour, SkipClean: true}
	out := offset.Apply(results, opts)
	if !out[0].DetectedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("drifted: expected shifted time, got %v", out[0].DetectedAt)
	}
	if !out[1].DetectedAt.Equal(now) {
		t.Errorf("clean: expected unchanged time, got %v", out[1].DetectedAt)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	original := []drift.Result{makeResults(true, now)}
	opts := offset.Options{Shift: time.Hour}
	offset.Apply(original, opts)
	if !original[0].DetectedAt.Equal(now) {
		t.Error("Apply mutated the original slice")
	}
}

func TestWrite_ForwardShift(t *testing.T) {
	var buf bytes.Buffer
	opts := offset.Options{Shift: time.Hour}
	offset.Write(&buf, opts, 3)
	got := buf.String()
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	for _, want := range []string{"3", "forward", "1h0m0s"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("expected output to contain %q, got: %s", want, got)
		}
	}
}

func TestWrite_BackwardShift_ShowsBackward(t *testing.T) {
	var buf bytes.Buffer
	opts := offset.Options{Shift: -2 * time.Hour}
	offset.Write(&buf, opts, 1)
	if !bytes.Contains(buf.Bytes(), []byte("backward")) {
		t.Errorf("expected 'backward' in output, got: %s", buf.String())
	}
}

package cascade_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/driftwatch/internal/cascade"
	"github.com/driftwatch/internal/drift"
)

func makeResults(n int, drifted bool) []drift.Result {
	out := make([]drift.Result, n)
	for i := range out {
		out[i] = drift.Result{
			Service: "svc",
			Field:   "image",
			Drifted: drifted,
		}
	}
	return out
}

func identity(r []drift.Result) ([]drift.Result, error) { return r, nil }

func halve(r []drift.Result) ([]drift.Result, error) {
	if len(r) == 0 {
		return r, nil
	}
	return r[:len(r)/2], nil
}

func failStage(_ []drift.Result) ([]drift.Result, error) {
	return nil, errors.New("stage failure")
}

func TestApply_NoStages_ReturnsOriginal(t *testing.T) {
	input := makeResults(4, true)
	out, errs := cascade.Apply(input, nil, cascade.DefaultOptions(), nil)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(out))
	}
}

func TestApply_SingleStage_TransformsResults(t *testing.T) {
	input := makeResults(8, false)
	stages := []cascade.Stage{{Name: "halve", Fn: halve}}
	out, errs := cascade.Apply(input, stages, cascade.DefaultOptions(), nil)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestApply_MultipleStages_ChainedCorrectly(t *testing.T) {
	input := makeResults(16, true)
	stages := []cascade.Stage{
		{Name: "first-halve", Fn: halve},
		{Name: "second-halve", Fn: halve},
	}
	out, errs := cascade.Apply(input, stages, cascade.DefaultOptions(), nil)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestApply_StopOnError_HaltsPipeline(t *testing.T) {
	input := makeResults(4, true)
	called := false
	after := func(r []drift.Result) ([]drift.Result, error) {
		called = true
		return r, nil
	}
	stages := []cascade.Stage{
		{Name: "fail", Fn: failStage},
		{Name: "after", Fn: after},
	}
	opts := cascade.DefaultOptions()
	opts.StopOnError = true
	_, errs := cascade.Apply(input, stages, opts, nil)
	if len(errs) == 0 {
		t.Fatal("expected an error")
	}
	if called {
		t.Fatal("pipeline should have stopped before the second stage")
	}
}

func TestApply_ContinueOnError_RunsRemainingStages(t *testing.T) {
	input := makeResults(4, true)
	called := false
	after := func(r []drift.Result) ([]drift.Result, error) {
		called = true
		return r, nil
	}
	stages := []cascade.Stage{
		{Name: "fail", Fn: failStage},
		{Name: "after", Fn: after},
	}
	opts := cascade.DefaultOptions()
	opts.StopOnError = false
	_, errs := cascade.Apply(input, stages, opts, nil)
	if len(errs) == 0 {
		t.Fatal("expected at least one error")
	}
	if !called {
		t.Fatal("subsequent stage should still have been called")
	}
}

func TestApply_Verbose_WritesStageInfo(t *testing.T) {
	input := makeResults(6, false)
	stages := []cascade.Stage{
		{Name: "identity", Fn: identity},
		{Name: "halve", Fn: halve},
	}
	opts := cascade.DefaultOptions()
	opts.Verbose = true
	var buf bytes.Buffer
	cascade.Apply(input, stages, opts, &buf)
	if !strings.Contains(buf.String(), "identity") {
		t.Error("expected stage name 'identity' in verbose output")
	}
	if !strings.Contains(buf.String(), "halve") {
		t.Error("expected stage name 'halve' in verbose output")
	}
}

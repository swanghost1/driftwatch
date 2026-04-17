package fanin_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/fanin"
)

func makeResult(service, field string, drifted bool) drift.Result {
	return drift.Result{Service: service, Field: field, Drifted: drifted}
}

func TestApply_NoSources_ReturnsNil(t *testing.T) {
	out := fanin.Apply(nil, fanin.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestApply_SingleSource_ReturnsAll(t *testing.T) {
	src := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-b", "replicas", false),
	}
	out := fanin.Apply([][]drift.Result{src}, fanin.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_MultipleSources_Combined(t *testing.T) {
	a := []drift.Result{makeResult("svc-a", "image", true)}
	b := []drift.Result{makeResult("svc-b", "replicas", false)}
	out := fanin.Apply([][]drift.Result{a, b}, fanin.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_Deduplication_RemovesDuplicates(t *testing.T) {
	a := []drift.Result{makeResult("svc-a", "image", true)}
	b := []drift.Result{makeResult("svc-a", "image", true)}
	out := fanin.Apply([][]drift.Result{a, b}, fanin.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 after dedup, got %d", len(out))
	}
}

func TestApply_DeduplicationDisabled_KeepsDuplicates(t *testing.T) {
	a := []drift.Result{makeResult("svc-a", "image", true)}
	b := []drift.Result{makeResult("svc-a", "image", true)}
	opts := fanin.Options{DriftedFirst: false, DeduplicateByKey: false}
	out := fanin.Apply([][]drift.Result{a, b}, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 without dedup, got %d", len(out))
	}
}

func TestApply_DriftedFirst_OrderCorrect(t *testing.T) {
	src := []drift.Result{
		makeResult("svc-a", "replicas", false),
		makeResult("svc-b", "image", true),
		makeResult("svc-c", "env", false),
		makeResult("svc-d", "replicas", true),
	}
	out := fanin.Apply([][]drift.Result{src}, fanin.DefaultOptions())
	for i, r := range out {
		if !r.Drifted && i < 2 {
			t.Errorf("expected drifted results first, got clean at index %d", i)
		}
	}
}

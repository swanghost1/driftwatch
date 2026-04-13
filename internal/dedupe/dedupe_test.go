package dedupe_test

import (
	"testing"

	"github.com/example/driftwatch/internal/dedupe"
	"github.com/example/driftwatch/internal/drift"
)

func makeResult(service, field, want, got string, drifted bool) drift.DetectResult {
	return drift.DetectResult{
		Service: service,
		Field:   field,
		Want:    want,
		Got:     got,
		Drifted: drifted,
	}
}

func TestApply_NoDuplicates_ReturnsAll(t *testing.T) {
	input := []drift.DetectResult{
		makeResult("svc-a", "image", "v1", "v1", false),
		makeResult("svc-b", "image", "v1", "v2", true),
	}
	got := dedupe.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_DuplicateRemoved(t *testing.T) {
	input := []drift.DetectResult{
		makeResult("svc-a", "image", "v1", "v2", true),
		makeResult("svc-a", "image", "v1", "v3", true), // duplicate key
	}
	got := dedupe.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 result after dedup, got %d", len(got))
	}
	if got[0].Got != "v2" {
		t.Errorf("expected first occurrence to be kept, got %q", got[0].Got)
	}
}

func TestApply_SameServiceDifferentFields_BothKept(t *testing.T) {
	input := []drift.DetectResult{
		makeResult("svc-a", "image", "v1", "v2", true),
		makeResult("svc-a", "replicas", "3", "1", true),
	}
	got := dedupe.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	got := dedupe.Apply(nil)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d", len(got))
	}
}

func TestApply_OKAndDriftSameKey_FirstKept(t *testing.T) {
	input := []drift.DetectResult{
		makeResult("svc-a", "image", "v1", "v1", false),
		makeResult("svc-a", "image", "v1", "v2", true),
	}
	got := dedupe.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Drifted {
		t.Error("expected first (non-drifted) entry to be kept")
	}
}

func TestCount_UniqueKeys(t *testing.T) {
	input := []drift.DetectResult{
		makeResult("svc-a", "image", "v1", "v2", true),
		makeResult("svc-a", "image", "v1", "v2", true), // duplicate
		makeResult("svc-b", "replicas", "2", "1", true),
	}
	if n := dedupe.Count(input); n != 2 {
		t.Errorf("expected count 2, got %d", n)
	}
}

func TestKey_String(t *testing.T) {
	k := dedupe.Key{Service: "my-svc", Field: "image"}
	if k.String() != "my-svc/image" {
		t.Errorf("unexpected key string: %q", k.String())
	}
}

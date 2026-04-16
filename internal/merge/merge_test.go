package merge_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/merge"
)

func makeResult(service, field, want, got string, drifted bool) drift.Result {
	return drift.Result{
		Service: service,
		Field:   field,
		Want:    want,
		Got:     got,
		Drifted: drifted,
	}
}

func TestApply_NoOverlap_ReturnsCombined(t *testing.T) {
	a := []drift.Result{makeResult("svc-a", "image", "v1", "v1", false)}
	b := []drift.Result{makeResult("svc-b", "image", "v2", "v3", true)}

	got := merge.Apply(a, b)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_DuplicateClean_KeepsFirst(t *testing.T) {
	a := []drift.Result{makeResult("svc-a", "image", "v1", "v1", false)}
	b := []drift.Result{makeResult("svc-a", "image", "v1", "v1", false)}

	got := merge.Apply(a, b)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestApply_ConflictDriftedWins(t *testing.T) {
	clean := []drift.Result{makeResult("svc-a", "image", "v1", "v1", false)}
	drifted := []drift.Result{makeResult("svc-a", "image", "v1", "v2", true)}

	got := merge.Apply(clean, drifted)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if !got[0].Drifted {
		t.Error("expected drifted result to win conflict")
	}
}

func TestApply_ConflictDriftedFirst_CleanSecond_DriftedWins(t *testing.T) {
	drifted := []drift.Result{makeResult("svc-a", "replicas", "3", "1", true)}
	clean := []drift.Result{makeResult("svc-a", "replicas", "3", "3", false)}

	got := merge.Apply(drifted, clean)
	if !got[0].Drifted {
		t.Error("drifted result should be retained when it appears first")
	}
}

func TestApply_SortedByServiceThenField(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-b", "image", "v1", "v1", false),
		makeResult("svc-a", "replicas", "2", "2", false),
		makeResult("svc-a", "image", "v1", "v1", false),
	}

	got := merge.Apply(input)
	if got[0].Service != "svc-a" || got[0].Field != "image" {
		t.Errorf("unexpected first result: %+v", got[0])
	}
	if got[1].Field != "replicas" {
		t.Errorf("unexpected second result: %+v", got[1])
	}
}

func TestApply_EmptyInputs_ReturnsEmpty(t *testing.T) {
	got := merge.Apply([]drift.Result{}, []drift.Result{})
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", "v1", "v2", true)}
	var buf bytes.Buffer
	merge.Write(&buf, results)
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header in output")
	}
}

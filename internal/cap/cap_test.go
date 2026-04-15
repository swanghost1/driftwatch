package cap_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/cap"
	"github.com/driftwatch/internal/drift"
)

func makeResult(service, field string, drifted bool) drift.Result {
	return drift.Result{
		Service: service,
		Field:   field,
		Drifted: drifted,
	}
}

func TestApply_ZeroMax_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-a", "replicas", true),
		makeResult("svc-b", "image", true),
	}
	opts := cap.Options{MaxPerService: 0}
	out := cap.Apply(results, opts)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_CapEnforced_LimitsPerService(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-a", "replicas", true),
		makeResult("svc-a", "env", true),
		makeResult("svc-b", "image", true),
	}
	opts := cap.Options{MaxPerService: 2}
	out := cap.Apply(results, opts)

	counts := map[string]int{}
	for _, r := range out {
		counts[r.Service]++
	}

	if counts["svc-a"] != 2 {
		t.Errorf("expected 2 results for svc-a, got %d", counts["svc-a"])
	}
	if counts["svc-b"] != 1 {
		t.Errorf("expected 1 result for svc-b, got %d", counts["svc-b"])
	}
}

func TestApply_CapOne_KeepsFirstResult(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-a", "replicas", false),
	}
	opts := cap.Options{MaxPerService: 1}
	out := cap.Apply(results, opts)

	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Field != "image" {
		t.Errorf("expected first field to be kept, got %q", out[0].Field)
	}
}

func TestApply_MultipleServices_IndependentCaps(t *testing.T) {
	results := []drift.Result{
		makeResult("alpha", "image", true),
		makeResult("beta", "image", true),
		makeResult("alpha", "replicas", true),
		makeResult("beta", "replicas", true),
		makeResult("gamma", "image", true),
	}
	opts := cap.Options{MaxPerService: 1}
	out := cap.Apply(results, opts)

	if len(out) != 3 {
		t.Errorf("expected 3 results (one per service), got %d", len(out))
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := cap.Apply(nil, cap.Options{MaxPerService: 5})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d", len(out))
	}
}

func TestWrite_ContainsSummaryFields(t *testing.T) {
	before := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-a", "replicas", true),
		makeResult("svc-b", "image", true),
	}
	opts := cap.Options{MaxPerService: 1}
	after := cap.Apply(before, opts)

	var buf bytes.Buffer
	cap.Write(&buf, before, after, opts)
	out := buf.String()

	for _, want := range []string{"max_per_service=1", "before=3", "after=2", "dropped=1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got: %s", want, out)
		}
	}
}

package reorder_test

import (
	"bytes"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/reorder"
)

func makeResult(service, field string, drifted bool, expected string) drift.Result {
	return drift.Result{
		Service:  service,
		Field:    field,
		Drifted:  drifted,
		Expected: expected,
	}
}

func TestApply_DefaultOrder_DriftedFirst(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", false, "nginx:1"),
		makeResult("svc-b", "image", true, "nginx:2"),
		makeResult("svc-c", "replicas", false, "3"),
	}
	out := reorder.Apply(input, reorder.DefaultOptions())
	if !out[0].Drifted {
		t.Fatalf("expected first result to be drifted, got service=%s", out[0].Service)
	}
}

func TestApply_ByService_Alphabetical(t *testing.T) {
	input := []drift.Result{
		makeResult("zebra", "image", false, ""),
		makeResult("alpha", "image", false, ""),
		makeResult("mango", "image", false, ""),
	}
	opts := reorder.Options{Order: []reorder.Criterion{reorder.ByService}}
	out := reorder.Apply(input, opts)
	if out[0].Service != "alpha" || out[1].Service != "mango" || out[2].Service != "zebra" {
		t.Fatalf("unexpected order: %v", serviceNames(out))
	}
}

func TestApply_ByField_Alphabetical(t *testing.T) {
	input := []drift.Result{
		makeResult("svc", "replicas", false, ""),
		makeResult("svc", "image", false, ""),
		makeResult("svc", "env", false, ""),
	}
	opts := reorder.Options{Order: []reorder.Criterion{reorder.ByField}}
	out := reorder.Apply(input, opts)
	if out[0].Field != "env" || out[1].Field != "image" || out[2].Field != "replicas" {
		t.Fatalf("unexpected field order: %v", fieldNames(out))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := []drift.Result{
		makeResult("z", "image", false, ""),
		makeResult("a", "image", false, ""),
	}
	opts := reorder.Options{Order: []reorder.Criterion{reorder.ByService}}
	_ = reorder.Apply(input, opts)
	if input[0].Service != "z" {
		t.Fatal("Apply mutated the original slice")
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := reorder.Apply(nil, reorder.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d results", len(out))
	}
}

func TestWrite_ContainsCriteria(t *testing.T) {
	var buf bytes.Buffer
	opts := reorder.Options{Order: []reorder.Criterion{reorder.ByDrifted, reorder.ByService}}
	reorder.Write(&buf, opts)
	out := buf.String()
	if !contains(out, "drifted") || !contains(out, "service") {
		t.Fatalf("Write output missing criteria: %q", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func serviceNames(rs []drift.Result) []string {
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = r.Service
	}
	return out
}

func fieldNames(rs []drift.Result) []string {
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = r.Field
	}
	return out
}

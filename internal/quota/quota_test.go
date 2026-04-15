package quota_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/quota"
)

func makeResults(pairs ...interface{}) []drift.Result {
	var out []drift.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		svc := pairs[i].(string)
		drifted := pairs[i+1].(bool)
		out = append(out, drift.Result{Service: svc, Drifted: drifted})
	}
	return out
}

func TestEvaluate_ZeroLimit_NoViolations(t *testing.T) {
	results := makeResults("api", true, "api", true, "api", true)
	r := quota.Evaluate(results, quota.Options{MaxDriftedFields: 0})
	if r.Exceeded() {
		t.Fatal("expected no violations when limit is zero")
	}
}

func TestEvaluate_NoDrift_NoViolations(t *testing.T) {
	results := makeResults("api", false, "worker", false)
	r := quota.Evaluate(results, quota.Options{MaxDriftedFields: 1})
	if r.Exceeded() {
		t.Fatalf("expected no violations, got %d", len(r.Violations))
	}
}

func TestEvaluate_WithinLimit_NoViolations(t *testing.T) {
	results := makeResults("api", true, "api", true)
	r := quota.Evaluate(results, quota.Options{MaxDriftedFields: 2})
	if r.Exceeded() {
		t.Fatalf("expected no violations for exactly-at-limit service")
	}
}

func TestEvaluate_ExceedsLimit_ReturnsViolation(t *testing.T) {
	results := makeResults("api", true, "api", true, "api", true)
	r := quota.Evaluate(results, quota.Options{MaxDriftedFields: 2})
	if !r.Exceeded() {
		t.Fatal("expected violation")
	}
	if len(r.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(r.Violations))
	}
	v := r.Violations[0]
	if v.Service != "api" {
		t.Errorf("expected service 'api', got %q", v.Service)
	}
	if v.DriftedFields != 3 {
		t.Errorf("expected 3 drifted fields, got %d", v.DriftedFields)
	}
	if v.Limit != 2 {
		t.Errorf("expected limit 2, got %d", v.Limit)
	}
}

func TestEvaluate_MultipleServices_OnlyViolatingReturned(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Drifted: true},
		{Service: "api", Drifted: true},
		{Service: "api", Drifted: true},
		{Service: "worker", Drifted: true},
	}
	r := quota.Evaluate(results, quota.Options{MaxDriftedFields: 2})
	if len(r.Violations) != 1 || r.Violations[0].Service != "api" {
		t.Errorf("expected only 'api' violation, got %+v", r.Violations)
	}
}

func TestWrite_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	err := quota.Write(&buf, quota.Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected 'no violations' in output, got: %s", buf.String())
	}
}

func TestWrite_WithViolations_ContainsServiceName(t *testing.T) {
	r := quota.Result{
		Violations: []quota.Violation{
			{Service: "api", DriftedFields: 5, Limit: 2},
		},
	}
	var buf bytes.Buffer
	if err := quota.Write(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"api", "5", "2", "SERVICE", "LIMIT"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

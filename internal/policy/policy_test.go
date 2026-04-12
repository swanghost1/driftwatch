package policy_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/policy"
)

func makeResults(driftedFields map[string][]string) []drift.Result {
	var results []drift.Result
	for svc, fields := range driftedFields {
		r := drift.Result{Service: svc}
		if len(fields) > 0 {
			r.Drifted = true
			for _, f := range fields {
				r.Diffs = append(r.Diffs, drift.Diff{Field: f, Want: "a", Got: "b"})
			}
		}
		results = append(results, r)
	}
	return results
}

func TestEvaluate_NoDrift_NoViolations(t *testing.T) {
	results := makeResults(map[string][]string{"svc-a": {}})
	rule := policy.Rule{MaxDriftedServices: 0}
	violations := policy.Evaluate(results, rule)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_DriftWithinLimit_NoViolations(t *testing.T) {
	results := makeResults(map[string][]string{
		"svc-a": {"replicas"},
	})
	rule := policy.Rule{MaxDriftedServices: 2}
	violations := policy.Evaluate(results, rule)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluate_ExceedsMaxDrifted_ReturnsViolation(t *testing.T) {
	results := makeResults(map[string][]string{
		"svc-a": {"replicas"},
		"svc-b": {"image"},
		"svc-c": {"image"},
	})
	rule := policy.Rule{MaxDriftedServices: 1}
	violations := policy.Evaluate(results, rule)
	if len(violations) == 0 {
		t.Fatal("expected at least one violation")
	}
	found := false
	for _, v := range violations {
		if v.Service == "" && v.Field == "" {
			found = true
		}
	}
	if !found {
		t.Error("expected a count-based violation with no service/field")
	}
}

func TestEvaluate_BlockedField_ReturnsViolation(t *testing.T) {
	results := makeResults(map[string][]string{
		"svc-a": {"image"},
	})
	rule := policy.Rule{MaxDriftedServices: 10, BlockedFields: []string{"image"}}
	violations := policy.Evaluate(results, rule)
	if len(violations) == 0 {
		t.Fatal("expected a blocked-field violation")
	}
	if violations[0].Service != "svc-a" || violations[0].Field != "image" {
		t.Errorf("unexpected violation: %+v", violations[0])
	}
}

func TestEvaluate_NonBlockedField_NoViolation(t *testing.T) {
	results := makeResults(map[string][]string{
		"svc-a": {"replicas"},
	})
	rule := policy.Rule{MaxDriftedServices: 10, BlockedFields: []string{"image"}}
	violations := policy.Evaluate(results, rule)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

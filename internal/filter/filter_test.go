package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/filter"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ServiceName: "api-gateway", Drifts: []drift.Item{{Field: "image", Expected: "v1", Actual: "v2"}}},
		{ServiceName: "auth-service", Drifts: []drift.Item{}},
		{ServiceName: "billing-worker", Drifts: []drift.Item{{Field: "replicas", Expected: "3", Actual: "2"}}},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{})
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_OnlyDrifted_ExcludesClean(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{OnlyDrifted: true})
	for _, r := range got {
		if len(r.Drifts) == 0 {
			t.Errorf("service %q has no drifts but was included", r.ServiceName)
		}
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(got))
	}
}

func TestApply_ServiceFilter_MatchesSubstring(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{Services: []string{"billing"}})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ServiceName != "billing-worker" {
		t.Errorf("unexpected service name %q", got[0].ServiceName)
	}
}

func TestApply_ServiceFilter_CaseInsensitive(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{Services: []string{"API"}})
	if len(got) != 1 || got[0].ServiceName != "api-gateway" {
		t.Fatalf("case-insensitive match failed, got %+v", got)
	}
}

func TestApply_ServiceFilter_MultiplePatterns(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{Services: []string{"auth", "billing"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	// auth-service matches name but has no drifts — should be excluded
	got := filter.Apply(makeResults(), filter.Options{
		Services:    []string{"auth", "billing"},
		OnlyDrifted: true,
	})
	if len(got) != 1 || got[0].ServiceName != "billing-worker" {
		t.Fatalf("expected only billing-worker, got %+v", got)
	}
}

func TestApply_NoMatchingService_ReturnsEmpty(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{Services: []string{"nonexistent"}})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

package groupfilter_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/groupfilter"
)

func makeResult(service string, drifted bool) drift.Result {
	r := drift.Result{Service: service}
	if drifted {
		r.Fields = []drift.FieldDrift{{Field: "image", Declared: "a", Live: "b"}}
	}
	return r
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("frontend/web", false),
		makeResult("backend/api", true),
	}
	got := groupfilter.Apply(results, groupfilter.Options{})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_IncludeGroups_FiltersOthers(t *testing.T) {
	results := []drift.Result{
		makeResult("frontend/web", false),
		makeResult("backend/api", true),
		makeResult("backend/worker", false),
	}
	got := groupfilter.Apply(results, groupfilter.Options{IncludeGroups: []string{"backend"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	for _, r := range got {
		if r.Service == "frontend/web" {
			t.Error("frontend/web should have been excluded")
		}
	}
}

func TestApply_ExcludeGroups_RemovesMatches(t *testing.T) {
	results := []drift.Result{
		makeResult("frontend/web", false),
		makeResult("backend/api", true),
	}
	got := groupfilter.Apply(results, groupfilter.Options{ExcludeGroups: []string{"frontend"}})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Service != "backend/api" {
		t.Errorf("unexpected service %q", got[0].Service)
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	results := []drift.Result{
		makeResult("Frontend/web", false),
		makeResult("backend/api", true),
	}
	got := groupfilter.Apply(results, groupfilter.Options{IncludeGroups: []string{"FRONTEND"}})
	if len(got) != 1 || got[0].Service != "Frontend/web" {
		t.Errorf("expected Frontend/web, got %v", got)
	}
}

func TestApply_NoSlash_UsesDefaultGroup(t *testing.T) {
	results := []drift.Result{
		makeResult("standalone", true),
		makeResult("backend/api", false),
	}
	got := groupfilter.Apply(results, groupfilter.Options{IncludeGroups: []string{"default"}})
	if len(got) != 1 || got[0].Service != "standalone" {
		t.Errorf("expected standalone, got %v", got)
	}
}

func TestGroups_ReturnsDistinct(t *testing.T) {
	results := []drift.Result{
		makeResult("frontend/web", false),
		makeResult("frontend/mobile", true),
		makeResult("backend/api", false),
		makeResult("standalone", true),
	}
	groups := groupfilter.Groups(results)
	if len(groups) != 3 {
		t.Errorf("expected 3 distinct groups, got %d: %v", len(groups), groups)
	}
}

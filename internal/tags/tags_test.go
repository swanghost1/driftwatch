package tags_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/tags"
)

func makeResult(name, env string, drifted bool) drift.Result {
	r := drift.Result{
		Service: name,
		Drifted: drifted,
		Tags:    map[string]string{},
	}
	if env != "" {
		r.Tags["env"] = env
	}
	return r
}

func TestGroupByTag_SplitsCorrectly(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "prod", false),
		makeResult("svc-b", "staging", true),
		makeResult("svc-c", "prod", true),
		makeResult("svc-d", "", false),
	}
	groups := tags.GroupByTag(results, "env")
	if len(groups["prod"]) != 2 {
		t.Errorf("expected 2 prod services, got %d", len(groups["prod"]))
	}
	if len(groups["staging"]) != 1 {
		t.Errorf("expected 1 staging service, got %d", len(groups["staging"]))
	}
	if len(groups["(untagged)"]) != 1 {
		t.Errorf("expected 1 untagged service, got %d", len(groups["(untagged)"]))
	}
}

func TestGroupByTag_AllUntagged(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "", false),
		makeResult("svc-b", "", true),
	}
	groups := tags.GroupByTag(results, "env")
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups["(untagged)"]) != 2 {
		t.Errorf("expected 2 untagged services, got %d", len(groups["(untagged)"]))
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "prod", false),
		makeResult("svc-b", "dev", false),
		makeResult("svc-c", "staging", false),
	}
	groups := tags.GroupByTag(results, "env")
	keys := tags.Keys(groups)
	expected := []string{"dev", "prod", "staging"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestFilterByTag_MatchesCaseInsensitive(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "Prod", false),
		makeResult("svc-b", "staging", true),
		makeResult("svc-c", "PROD", true),
	}
	out := tags.FilterByTag(results, "env", "prod")
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestFilterByTag_EmptyVal_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "prod", false),
		makeResult("svc-b", "staging", true),
	}
	out := tags.FilterByTag(results, "env", "")
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

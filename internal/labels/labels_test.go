package labels_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/labels"
)

func makeResult(name string, lbls map[string]string) drift.Result {
	return drift.Result{
		Service: name,
		Labels:  lbls,
	}
}

func TestParseLabel_Valid(t *testing.T) {
	k, v, err := labels.ParseLabel("env=production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k != "env" || v != "production" {
		t.Errorf("got key=%q value=%q, want env/production", k, v)
	}
}

func TestParseLabel_InvalidFormat(t *testing.T) {
	for _, raw := range []string{"nodivider", "=value", ""} {
		_, _, err := labels.ParseLabel(raw)
		if err == nil {
			t.Errorf("expected error for %q, got nil", raw)
		}
	}
}

func TestParseAll_ReturnsMap(t *testing.T) {
	m, err := labels.ParseAll([]string{"team=platform", "env=staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["team"] != "platform" || m["env"] != "staging" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestParseAll_PropagatesError(t *testing.T) {
	_, err := labels.ParseAll([]string{"valid=yes", "bad"})
	if err == nil {
		t.Fatal("expected error for malformed label")
	}
}

func TestMatch_AllPresent_ReturnsTrue(t *testing.T) {
	r := makeResult("svc", map[string]string{"env": "prod", "team": "platform"})
	if !labels.Match(r, map[string]string{"env": "prod"}) {
		t.Error("expected match")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	r := makeResult("svc", map[string]string{"ENV": "PROD"})
	if !labels.Match(r, map[string]string{"env": "prod"}) {
		t.Error("expected case-insensitive match")
	}
}

func TestMatch_MissingKey_ReturnsFalse(t *testing.T) {
	r := makeResult("svc", map[string]string{"team": "platform"})
	if labels.Match(r, map[string]string{"env": "prod"}) {
		t.Error("expected no match")
	}
}

func TestFilter_ReturnsMatchingResults(t *testing.T) {
	results := []drift.Result{
		makeResult("a", map[string]string{"env": "prod"}),
		makeResult("b", map[string]string{"env": "staging"}),
		makeResult("c", map[string]string{"env": "prod", "team": "infra"}),
	}
	got := labels.Filter(results, map[string]string{"env": "prod"})
	if len(got) != 2 {
		t.Errorf("expected 2 results, got %d", len(got))
	}
}

func TestFilter_EmptyRequired_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("a", nil),
		makeResult("b", nil),
	}
	got := labels.Filter(results, map[string]string{})
	if len(got) != 2 {
		t.Errorf("expected 2 results, got %d", len(got))
	}
}

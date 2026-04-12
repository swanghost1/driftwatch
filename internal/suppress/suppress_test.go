package suppress_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/suppress"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ServiceName: "api", Field: "image", Expected: "v1", Actual: "v2", Drifted: true},
		{ServiceName: "api", Field: "replicas", Expected: "3", Actual: "2", Drifted: true},
		{ServiceName: "worker", Field: "image", Expected: "v1", Actual: "v1", Drifted: false},
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	results := makeResults()
	out := suppress.Apply(results, &suppress.Store{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_MatchingRule_SuppressesResult(t *testing.T) {
	results := makeResults()
	s := &suppress.Store{
		Rules: []suppress.Rule{
			{Service: "api", Field: "image", Reason: "known drift"},
		},
	}
	out := suppress.Apply(results, s)
	if len(out) != 2 {
		t.Fatalf("expected 2 results after suppression, got %d", len(out))
	}
	for _, r := range out {
		if r.ServiceName == "api" && r.Field == "image" {
			t.Error("suppressed result should not appear in output")
		}
	}
}

func TestApply_ExpiredRule_DoesNotSuppress(t *testing.T) {
	results := makeResults()
	s := &suppress.Store{
		Rules: []suppress.Rule{
			{Service: "api", Field: "image", Reason: "old", Expires: time.Now().Add(-time.Hour)},
		},
	}
	out := suppress.Apply(results, s)
	if len(out) != len(results) {
		t.Fatalf("expected all %d results, got %d", len(results), len(out))
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	results := makeResults()
	s := &suppress.Store{
		Rules: []suppress.Rule{
			{Service: "API", Field: "IMAGE", Reason: "case test"},
		},
	}
	out := suppress.Apply(results, s)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppress.json")

	orig := &suppress.Store{
		Rules: []suppress.Rule{
			{Service: "svc", Field: "replicas", Reason: "scaling test"},
		},
	}
	if err := suppress.SaveStore(path, orig); err != nil {
		t.Fatalf("SaveStore: %v", err)
	}
	loaded, err := suppress.LoadStore(path)
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].Service != "svc" {
		t.Errorf("unexpected loaded rules: %+v", loaded.Rules)
	}
}

func TestLoadStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	s, err := suppress.LoadStore("/nonexistent/path/suppress.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.Rules) != 0 {
		t.Errorf("expected empty store, got %+v", s)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "suppress.json")
	s := &suppress.Store{Rules: []suppress.Rule{}}
	if err := suppress.SaveStore(path, s); err != nil {
		t.Fatalf("SaveStore: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
	var loaded suppress.Store
	data, _ := os.ReadFile(path)
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Errorf("invalid JSON written: %v", err)
	}
}

package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ServiceName: "api", Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.25", Drifted: true},
		{ServiceName: "worker", Field: "replicas", Expected: "3", Actual: "2", Drifted: true},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := makeResults()

	if err := baseline.Save(path, "v1", results); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	entry, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if entry.Label != "v1" {
		t.Errorf("expected label v1, got %q", entry.Label)
	}
	if len(entry.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(entry.Results))
	}
	if entry.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "baseline.json")

	if err := baseline.Save(path, "test", makeResults()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %s: %v", path, err)
	}
}

func TestCompare_ReturnsOnlyNovelDrifts(t *testing.T) {
	existing := &baseline.Entry{
		CreatedAt: time.Now(),
		Label:     "old",
		Results: []drift.Result{
			{ServiceName: "api", Field: "image", Drifted: true},
		},
	}

	current := []drift.Result{
		{ServiceName: "api", Field: "image", Drifted: true},   // already in baseline
		{ServiceName: "api", Field: "replicas", Drifted: true}, // new
	}

	novel := baseline.Compare(existing, current)
	if len(novel) != 1 {
		t.Fatalf("expected 1 novel result, got %d", len(novel))
	}
	if novel[0].Field != "replicas" {
		t.Errorf("expected novel field 'replicas', got %q", novel[0].Field)
	}
}

func TestCompare_EmptyBaseline_ReturnsAll(t *testing.T) {
	existing := &baseline.Entry{Results: []drift.Result{}}
	current := makeResults()

	novel := baseline.Compare(existing, current)
	if len(novel) != len(current) {
		t.Errorf("expected %d results, got %d", len(current), len(novel))
	}
}

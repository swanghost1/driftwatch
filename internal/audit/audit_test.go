package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftwatch/internal/audit"
	"github.com/example/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ServiceName: "api", Drifted: false},
		{ServiceName: "worker", Drifted: true, Diffs: []drift.Diff{
			{Field: "image", Expected: "worker:v1", Actual: "worker:v2"},
		}},
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	l := audit.NewLogger(dir)

	if err := l.Record("driftwatch.yaml", "cli", makeResults()); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	entries, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 audit file, got %d", len(entries))
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "audit")
	l := audit.NewLogger(dir)

	if err := l.Record("cfg.yaml", "schedule", makeResults()); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestList_ReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	l := audit.NewLogger(dir)

	results := makeResults()
	if err := l.Record("driftwatch.yaml", "cli", results); err != nil {
		t.Fatalf("Record() error: %v", err)
	}

	entries, err := l.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	e := entries[0]
	if e.ConfigFile != "driftwatch.yaml" {
		t.Errorf("ConfigFile = %q, want %q", e.ConfigFile, "driftwatch.yaml")
	}
	if e.TriggeredBy != "cli" {
		t.Errorf("TriggeredBy = %q, want %q", e.TriggeredBy, "cli")
	}
	if e.TotalChecked != 2 {
		t.Errorf("TotalChecked = %d, want 2", e.TotalChecked)
	}
	if e.DriftedCount != 1 {
		t.Errorf("DriftedCount = %d, want 1", e.DriftedCount)
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	l := audit.NewLogger(dir)

	entries, err := l.List()
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}
}

func TestRecord_DriftedCountAccurate(t *testing.T) {
	dir := t.TempDir()
	l := audit.NewLogger(dir)

	results := []drift.Result{
		{ServiceName: "a", Drifted: true},
		{ServiceName: "b", Drifted: true},
		{ServiceName: "c", Drifted: false},
	}
	_ = l.Record("cfg.yaml", "test", results)

	entries, _ := l.List()
	if entries[0].DriftedCount != 2 {
		t.Errorf("DriftedCount = %d, want 2", entries[0].DriftedCount)
	}
}

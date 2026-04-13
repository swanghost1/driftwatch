package variance_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/variance"
)

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	store := variance.NewStore(dir)

	results := []drift.Result{
		{Service: "svc-a", Drifted: true, Field: "image"},
	}
	if err := store.Record(results); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "a", "b", "c")
	store := variance.NewStore(dir)

	if err := store.Record(nil); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}

func TestLoad_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	store := variance.NewStore(dir)

	entries, err := store.Load(0)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}
}

func TestLoad_MissingDir_ReturnsNil(t *testing.T) {
	store := variance.NewStore("/tmp/driftwatch-variance-nonexistent-xyz")
	entries, err := store.Load(0)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil for missing dir")
	}
}

func TestRoundTrip_PreservesResults(t *testing.T) {
	dir := t.TempDir()
	store := variance.NewStore(dir)

	want := []drift.Result{
		{Service: "svc-a", Drifted: true, Field: "replicas"},
		{Service: "svc-b", Drifted: false},
	}
	if err := store.Record(want); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := store.Load(0)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(entries[0].Results))
	}
}

func TestLoad_LimitRespected(t *testing.T) {
	dir := t.TempDir()
	store := variance.NewStore(dir)

	for i := 0; i < 5; i++ {
		if err := store.Record(nil); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	entries, err := store.Load(3)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries with limit=3, got %d", len(entries))
	}
}

package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/checkpoint"
)

func makeEntry(name string) checkpoint.Entry {
	return checkpoint.Entry{
		Name:    name,
		RunAt:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Total:   10,
		Drifted: 3,
		Clean:   7,
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := checkpoint.NewStore(dir)
	e := makeEntry("production")

	if err := store.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load("production")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != e.Name {
		t.Errorf("Name: got %q, want %q", got.Name, e.Name)
	}
	if got.Drifted != e.Drifted {
		t.Errorf("Drifted: got %d, want %d", got.Drifted, e.Drifted)
	}
	if !got.RunAt.Equal(e.RunAt) {
		t.Errorf("RunAt: got %v, want %v", got.RunAt, e.RunAt)
	}
}

func TestLoad_FileNotFound_ReturnsErrNotFound(t *testing.T) {
	store := checkpoint.NewStore(t.TempDir())
	_, err := store.Load("nonexistent")
	if err != checkpoint.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "a", "b", "c")
	store := checkpoint.NewStore(dir)

	if err := store.Save(makeEntry("svc")); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "svc.json")); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestSave_Overwrites_PreviousEntry(t *testing.T) {
	store := checkpoint.NewStore(t.TempDir())

	e1 := makeEntry("staging")
	e1.Drifted = 1
	if err := store.Save(e1); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	e2 := makeEntry("staging")
	e2.Drifted = 9
	if err := store.Save(e2); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	got, err := store.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Drifted != 9 {
		t.Errorf("Drifted: got %d, want 9", got.Drifted)
	}
}

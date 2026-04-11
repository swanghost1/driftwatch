package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/snapshot"
)

func makeStore() *snapshot.Store {
	return &snapshot.Store{
		Version: "1",
		Snapshots: []snapshot.ServiceSnapshot{
			{
				Name:       "api",
				Image:      "api:v1.2.3",
				Replicas:   3,
				Env:        map[string]string{"LOG_LEVEL": "info"},
				CapturedAt: time.Now().UTC(),
			},
			{
				Name:       "worker",
				Image:      "worker:v0.9.0",
				Replicas:   1,
				CapturedAt: time.Now().UTC(),
			},
		},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := makeStore()
	if err := snapshot.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Version != orig.Version {
		t.Errorf("version: got %q, want %q", loaded.Version, orig.Version)
	}
	if len(loaded.Snapshots) != len(orig.Snapshots) {
		t.Fatalf("snapshot count: got %d, want %d", len(loaded.Snapshots), len(orig.Snapshots))
	}
	if loaded.Snapshots[0].Image != "api:v1.2.3" {
		t.Errorf("image: got %q", loaded.Snapshots[0].Image)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "nested", "snap.json")

	if err := snapshot.Save(path, makeStore()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestFindByName(t *testing.T) {
	store := makeStore()

	got := store.FindByName("api")
	if got == nil || got.Image != "api:v1.2.3" {
		t.Errorf("FindByName(api): unexpected result %v", got)
	}

	missing := store.FindByName("unknown")
	if missing != nil {
		t.Errorf("expected nil for unknown service, got %v", missing)
	}
}

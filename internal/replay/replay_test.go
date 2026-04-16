package replay_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/replay"
	"github.com/example/driftwatch/internal/snapshot"
)

func makeStore(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(path, snapshot.Store{Entries: entries}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return path
}

func baseEntry(name string) snapshot.Entry {
	return snapshot.Entry{
		Name:        name,
		Image:       "nginx:1.25",
		Replicas:    2,
		CapturedAt:  time.Now().UTC(),
	}
}

func TestRun_NoDrift(t *testing.T) {
	e := baseEntry("web")
	path := makeStore(t, []snapshot.Entry{e})

	results, err := replay.Run([]snapshot.Entry{e}, replay.Options{SnapshotPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Error("expected no drift")
	}
}

func TestRun_ImageDrift(t *testing.T) {
	live := baseEntry("web")
	decl := live
	decl.Image = "nginx:1.26"
	path := makeStore(t, []snapshot.Entry{live})

	results, err := replay.Run([]snapshot.Entry{decl}, replay.Options{SnapshotPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || !results[0].Drifted {
		t.Error("expected drift on image field")
	}
}

func TestRun_AsOf_ExcludesFutureSnapshots(t *testing.T) {
	e := baseEntry("web")
	e.CapturedAt = time.Now().Add(time.Hour)
	path := makeStore(t, []snapshot.Entry{e})

	results, err := replay.Run([]snapshot.Entry{e}, replay.Options{
		SnapshotPath: path,
		AsOf:         time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRun_MissingSnapshotFile_ReturnsError(t *testing.T) {
	_, err := replay.Run(nil, replay.Options{SnapshotPath: "/no/such/file.json"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_ServiceNotInSnapshot_Skipped(t *testing.T) {
	e := baseEntry("api")
	path := makeStore(t, []snapshot.Entry{e})

	decl := baseEntry("web") // different name
	results, err := replay.Run([]snapshot.Entry{decl}, replay.Options{SnapshotPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRun_ReplayedAt_IsSet(t *testing.T) {
	e := baseEntry("svc")
	path := makeStore(t, []snapshot.Entry{e})
	before := time.Now().UTC()

	results, err := replay.Run([]snapshot.Entry{e}, replay.Options{SnapshotPath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].ReplayedAt.Before(before) {
		t.Error("ReplayedAt should be set to approximately now")
	}
}

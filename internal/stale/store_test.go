package stale_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/stale"
)

func TestStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stale.json")
	s := stale.NewStore(path)

	in := map[string]time.Time{
		"svc-a:image":    epoch,
		"svc-b:replicas": epoch.Add(-24 * time.Hour),
	}
	if err := s.Save(in); err != nil {
		t.Fatalf("Save: %v", err)
	}
	out, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(out) != len(in) {
		t.Fatalf("expected %d entries, got %d", len(in), len(out))
	}
	for k, v := range in {
		if !out[k].Equal(v) {
			t.Errorf("key %s: expected %v, got %v", k, v, out[k])
		}
	}
}

func TestStore_Load_FileNotFound_ReturnsEmptyMap(t *testing.T) {
	dir := t.TempDir()
	s := stale.NewStore(filepath.Join(dir, "missing.json"))
	m, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil || len(m) != 0 {
		t.Error("expected empty non-nil map")
	}
}

func TestStore_Save_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "stale.json")
	s := stale.NewStore(path)
	if err := s.Save(map[string]time.Time{"k": epoch}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	out, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

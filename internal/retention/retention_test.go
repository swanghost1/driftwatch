package retention

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeEntry(t *testing.T, dir string, ts time.Time, name string) string {
	t.Helper()
	p := filepath.Join(dir, name+".json")
	if err := os.WriteFile(p, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("writeEntry: %v", err)
	}
	if err := os.Chtimes(p, ts, ts); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return p
}

func TestPrune_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	pruned, err := Prune(dir, DefaultPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pruned) != 0 {
		t.Errorf("expected no pruned files, got %d", len(pruned))
	}
}

func TestPrune_OldEntries_Removed(t *testing.T) {
	dir := t.TempDir()
	old := time.Now().Add(-60 * 24 * time.Hour)
	writeEntry(t, dir, old, "old_entry")
	writeEntry(t, dir, time.Now(), "new_entry")

	p := Policy{MaxAge: 30 * 24 * time.Hour, MaxEntries: 0}
	pruned, err := Prune(dir, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pruned) != 1 {
		t.Errorf("expected 1 pruned file, got %d", len(pruned))
	}
}

func TestPrune_MaxEntries_KeepsNewest(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 5; i++ {
		ts := time.Now().Add(time.Duration(-i) * time.Hour)
		writeEntry(t, dir, ts, filepath.Base(ts.Format(time.RFC3339)))
	}

	p := Policy{MaxAge: 0, MaxEntries: 3}
	pruned, err := Prune(dir, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pruned) != 2 {
		t.Errorf("expected 2 pruned files, got %d", len(pruned))
	}

	remaining, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if len(remaining) != 3 {
		t.Errorf("expected 3 remaining files, got %d", len(remaining))
	}
}

func TestPrune_WithinLimits_NoPrune(t *testing.T) {
	dir := t.TempDir()
	writeEntry(t, dir, time.Now().Add(-1*time.Hour), "recent")

	p := Policy{MaxAge: 24 * time.Hour, MaxEntries: 10}
	pruned, err := Prune(dir, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pruned) != 0 {
		t.Errorf("expected no pruned files, got %d", len(pruned))
	}
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := DefaultPolicy()
	if p.MaxAge != 30*24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", p.MaxAge)
	}
	if p.MaxEntries != 100 {
		t.Errorf("unexpected MaxEntries: %d", p.MaxEntries)
	}
}

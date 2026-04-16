package spike

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "spike.json"))
	if err := s.Record(3, time.Now()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "spike.json")); err != nil {
		t.Fatal("expected file to exist")
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "a", "b", "spike.json"))
	if err := s.Record(1, time.Now()); err != nil {
		t.Fatalf("Record: %v", err)
	}
}

func TestCounts_EmptyStore_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "spike.json"))
	counts, err := s.Counts()
	if err != nil {
		t.Fatalf("Counts: %v", err)
	}
	if len(counts) != 0 {
		t.Fatalf("expected empty, got %v", counts)
	}
}

func TestCounts_ReturnsSortedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "spike.json"))
	now := time.Now()
	_ = s.Record(10, now)
	_ = s.Record(2, now.Add(-2*time.Hour))
	_ = s.Record(5, now.Add(-1*time.Hour))

	counts, err := s.Counts()
	if err != nil {
		t.Fatalf("Counts: %v", err)
	}
	if len(counts) != 3 {
		t.Fatalf("expected 3 counts, got %d", len(counts))
	}
	if counts[0] != 2 || counts[1] != 5 || counts[2] != 10 {
		t.Fatalf("unexpected order: %v", counts)
	}
}

func TestRoundTrip_PreservesCounts(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "spike.json"))
	now := time.Now()
	for i, v := range []int{1, 3, 7} {
		_ = s.Record(v, now.Add(time.Duration(i)*time.Minute))
	}
	counts, _ := s.Counts()
	if len(counts) != 3 {
		t.Fatalf("expected 3, got %d", len(counts))
	}
}

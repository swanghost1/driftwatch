package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/history"
)

func makeResults(drifted bool) []drift.Result {
	status := drift.StatusOK
	if drifted {
		status = drift.StatusDrifted
	}
	return []drift.Result{
		{ServiceName: "api", Status: status, Diffs: nil},
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	s := history.NewStore(dir)
	if err := s.Record(makeResults(false)); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestList_SortedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	s := history.NewStore(dir)
	for i := 0; i < 3; i++ {
		if err := s.Record(makeResults(i%2 == 0)); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
		time.Sleep(time.Millisecond * 10)
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].RunAt.Before(entries[i-1].RunAt) {
			t.Errorf("entries not sorted: %v before %v", entries[i].RunAt, entries[i-1].RunAt)
		}
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	s := history.NewStore(dir)
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_CreatesIntermediateDirectories(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/history"
	s := history.NewStore(dir)
	if err := s.Record(makeResults(true)); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}

func TestEntry_ResultsPreserved(t *testing.T) {
	dir := t.TempDir()
	s := history.NewStore(dir)
	results := makeResults(true)
	if err := s.Record(results); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(entries[0].Results))
	}
	if entries[0].Results[0].ServiceName != "api" {
		t.Errorf("service name mismatch: %s", entries[0].Results[0].ServiceName)
	}
}

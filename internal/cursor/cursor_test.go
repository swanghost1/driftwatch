package cursor_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/cursor"
)

func newStore(t *testing.T) *cursor.Store {
	t.Helper()
	return cursor.NewStore(t.TempDir())
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := newStore(t)
	st := cursor.State{Name: "svc-a", Offset: 42, RunCount: 3, LastSeen: time.Now().UTC().Truncate(time.Second)}
	if err := s.Save(st); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("svc-a")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Offset != st.Offset || got.RunCount != st.RunCount || got.Name != st.Name {
		t.Errorf("got %+v, want %+v", got, st)
	}
}

func TestLoad_FileNotFound_ReturnsErrNotFound(t *testing.T) {
	s := newStore(t)
	_, err := s.Load("missing")
	if err != cursor.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	base := filepath.Join(t.TempDir(), "a", "b", "c")
	s := cursor.NewStore(base)
	st := cursor.State{Name: "x", Offset: 1}
	if err := s.Save(st); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(base); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestAdvance_IncrementsOffsetAndRunCount(t *testing.T) {
	st := cursor.State{Name: "svc", Offset: 10, RunCount: 2}
	out := cursor.Advance(st, 5)
	if out.Offset != 15 {
		t.Errorf("Offset: got %d, want 15", out.Offset)
	}
	if out.RunCount != 3 {
		t.Errorf("RunCount: got %d, want 3", out.RunCount)
	}
	if out.LastSeen.IsZero() {
		t.Error("LastSeen should be set")
	}
}

func TestAdvance_DoesNotMutateOriginal(t *testing.T) {
	st := cursor.State{Name: "svc", Offset: 5, RunCount: 1}
	_ = cursor.Advance(st, 3)
	if st.Offset != 5 {
		t.Error("original state was mutated")
	}
}

func TestSave_Overwrites_PreviousState(t *testing.T) {
	s := newStore(t)
	s.Save(cursor.State{Name: "svc", Offset: 1})
	s.Save(cursor.State{Name: "svc", Offset: 99})
	got, _ := s.Load("svc")
	if got.Offset != 99 {
		t.Errorf("expected 99, got %d", got.Offset)
	}
}

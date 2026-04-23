package watermark_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/watermark"
)

func drifted(service string) drift.Result {
	return drift.Result{Service: service, Drifted: true}
}

func clean(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestLoad_FileNotFound_ReturnsZeroMark(t *testing.T) {
	s := watermark.NewStore(filepath.Join(t.TempDir(), "mark.json"))
	m, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Peak != 0 {
		t.Errorf("expected peak 0, got %d", m.Peak)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := watermark.NewStore(filepath.Join(t.TempDir(), "mark.json"))
	want := watermark.Mark{Peak: 7, RecordedAt: time.Now().UTC().Truncate(time.Second)}
	if err := s.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Peak != want.Peak {
		t.Errorf("peak: got %d, want %d", got.Peak, want.Peak)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c", "mark.json")
	s := watermark.NewStore(path)
	if err := s.Save(watermark.Mark{Peak: 1, RecordedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestUpdate_NewPeak_ReturnsTrueAndUpdates(t *testing.T) {
	s := watermark.NewStore(filepath.Join(t.TempDir(), "mark.json"))
	results := []drift.Result{drifted("svc-a"), drifted("svc-b"), clean("svc-c")}
	m, updated, err := s.Update(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Error("expected updated=true")
	}
	if m.Peak != 2 {
		t.Errorf("peak: got %d, want 2", m.Peak)
	}
}

func TestUpdate_BelowPeak_ReturnsFalse(t *testing.T) {
	s := watermark.NewStore(filepath.Join(t.TempDir(), "mark.json"))
	// Establish a peak of 3.
	_ = s.Save(watermark.Mark{Peak: 3, RecordedAt: time.Now().UTC()})
	results := []drift.Result{drifted("svc-a"), clean("svc-b")}
	m, updated, err := s.Update(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated {
		t.Error("expected updated=false")
	}
	if m.Peak != 3 {
		t.Errorf("peak should remain 3, got %d", m.Peak)
	}
}

func TestWrite_NoPeak_ShowsNoDriftMessage(t *testing.T) {
	var buf bytes.Buffer
	watermark.Write(&buf, watermark.Mark{})
	if got := buf.String(); got != "high-water mark: no drift recorded\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestWrite_WithPeak_ContainsPeakCount(t *testing.T) {
	var buf bytes.Buffer
	watermark.Write(&buf, watermark.Mark{Peak: 5, RecordedAt: time.Now().UTC()})
	if got := buf.String(); len(got) == 0 {
		t.Error("expected non-empty output")
	}
	const want = "5 drifted service(s)"
	if !bytes.Contains([]byte(got), []byte(want)) {
		t.Errorf("output %q does not contain %q", got, want)
	}
}

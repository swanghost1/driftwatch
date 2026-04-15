package window_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/window"
)

func TestNewSummary_Counts(t *testing.T) {
	now := time.Now()
	original := []drift.Result{
		makeResult("a", now.Add(-1*time.Hour), true),
		makeResult("b", now.Add(-3*time.Hour), false),
	}
	filtered := []drift.Result{original[0]}
	s := window.NewSummary(original, filtered, window.Options{Since: now.Add(-2*time.Hour)})
	if s.Matched != 1 {
		t.Errorf("expected Matched=1, got %d", s.Matched)
	}
	if s.Filtered != 1 {
		t.Errorf("expected Filtered=1, got %d", s.Filtered)
	}
	if s.Since == nil {
		t.Error("expected Since to be set")
	}
	if s.Until != nil {
		t.Error("expected Until to be nil")
	}
}

func TestWriteText_ContainsHeaders(t *testing.T) {
	s := window.Summary{Matched: 3, Filtered: 1}
	var buf bytes.Buffer
	if err := window.WriteText(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Window", "Matched", "Filtered", "(open)"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteText_ShowsBounds(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	s := window.Summary{Since: &now, Matched: 2, Filtered: 0}
	var buf bytes.Buffer
	_ = window.WriteText(&buf, s)
	if !strings.Contains(buf.String(), now.Format(time.RFC3339)) {
		t.Errorf("expected time in output")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	s := window.Summary{Matched: 5, Filtered: 2}
	var buf bytes.Buffer
	if err := window.WriteJSON(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if v, ok := out["matched"]; !ok || int(v.(float64)) != 5 {
		t.Errorf("expected matched=5")
	}
}

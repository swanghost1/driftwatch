package replay_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/replay"
)

func makeReplayResults() []replay.Result {
	return []replay.Result{
		{
			Result:     drift.Result{Service: "web", Field: "image", Drifted: true, Want: "nginx:1.26", Got: "nginx:1.25"},
			ReplayedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			Result:     drift.Result{Service: "api", Field: "replicas", Drifted: false, Want: "3", Got: "3"},
			ReplayedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		},
	}
}

func TestWriteText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := replay.WriteText(&buf, makeReplayResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header row with SERVICE")
	}
}

func TestWriteText_DriftedRow_ShowsYES(t *testing.T) {
	var buf bytes.Buffer
	if err := replay.WriteText(&buf, makeReplayResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "YES") {
		t.Error("expected YES for drifted row")
	}
}

func TestWriteText_CleanRow_ShowsNo(t *testing.T) {
	var buf bytes.Buffer
	if err := replay.WriteText(&buf, makeReplayResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no") {
		t.Error("expected 'no' for clean row")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := replay.WriteJSON(&buf, makeReplayResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestWriteJSON_ContainsReplayedAt(t *testing.T) {
	var buf bytes.Buffer
	if err := replay.WriteJSON(&buf, makeReplayResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ReplayedAt") {
		t.Error("expected ReplayedAt field in JSON output")
	}
}

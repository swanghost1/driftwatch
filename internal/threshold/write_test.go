package threshold_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/threshold"
)

func TestWriteJSON_ValidJSON(t *testing.T) {
	r := threshold.Result{Total: 10, Drifted: 4, DriftRate: 0.4, Breached: false}
	var sb strings.Builder
	if err := threshold.WriteJSON(&sb, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(sb.String()), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["total"].(float64) != 10 {
		t.Errorf("unexpected total: %v", m["total"])
	}
}

func TestWriteJSON_BreachedContainsReason(t *testing.T) {
	r := threshold.Result{Total: 5, Drifted: 5, DriftRate: 1.0, Breached: true, Reason: "count threshold"}
	var sb strings.Builder
	_ = threshold.WriteJSON(&sb, r)
	if !strings.Contains(sb.String(), "count threshold") {
		t.Error("expected reason in JSON output")
	}
}

func TestWriteText_OKLine(t *testing.T) {
	r := threshold.Result{Total: 4, Drifted: 1, DriftRate: 0.25}
	var sb strings.Builder
	threshold.WriteText(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "threshold=ok") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "drifted=1/4") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestWriteText_BreachedLine(t *testing.T) {
	r := threshold.Result{Total: 4, Drifted: 4, DriftRate: 1.0, Breached: true, Reason: "rate"}
	var sb strings.Builder
	threshold.WriteText(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "threshold=breached") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "reason") {
		t.Errorf("expected reason in text output: %s", out)
	}
}

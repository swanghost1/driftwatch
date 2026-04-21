package headroom_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/headroom"
)

func makeResults(driftedCount, cleanCount int) []drift.Result {
	var results []drift.Result
	for i := 0; i < driftedCount; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < cleanCount; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: false})
	}
	return results
}

func TestCompute_NoDrift(t *testing.T) {
	r := headroom.Compute(makeResults(0, 5), headroom.DefaultOptions())
	if r.Drifted != 0 {
		t.Fatalf("expected 0 drifted, got %d", r.Drifted)
	}
	if r.DriftRatePct != 0 {
		t.Fatalf("expected 0%% rate, got %.2f", r.DriftRatePct)
	}
}

func TestCompute_AllDrifted(t *testing.T) {
	r := headroom.Compute(makeResults(4, 0), headroom.DefaultOptions())
	if r.DriftRatePct != 100 {
		t.Fatalf("expected 100%% rate, got %.2f", r.DriftRatePct)
	}
}

func TestCompute_CountHeadroom_NoLimit(t *testing.T) {
	r := headroom.Compute(makeResults(2, 3), headroom.DefaultOptions())
	if r.CountHeadroom != -1 {
		t.Fatalf("expected -1 (no limit), got %d", r.CountHeadroom)
	}
}

func TestCompute_CountHeadroom_WithLimit(t *testing.T) {
	opts := headroom.Options{MaxDrifted: 5}
	r := headroom.Compute(makeResults(3, 2), opts)
	if r.CountHeadroom != 2 {
		t.Fatalf("expected headroom 2, got %d", r.CountHeadroom)
	}
}

func TestCompute_RateHeadroom_WithLimit(t *testing.T) {
	opts := headroom.Options{MaxRatePct: 50}
	// 2 drifted out of 4 = 50%, headroom = 0
	r := headroom.Compute(makeResults(2, 2), opts)
	if r.RateHeadroom != 0 {
		t.Fatalf("expected rate headroom 0, got %.2f", r.RateHeadroom)
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	r := headroom.Compute(nil, headroom.DefaultOptions())
	if r.Total != 0 || r.DriftRatePct != 0 {
		t.Fatal("expected zero report for empty input")
	}
}

func TestWriteText_ContainsKeyFields(t *testing.T) {
	opts := headroom.Options{MaxDrifted: 10, MaxRatePct: 80}
	r := headroom.Compute(makeResults(3, 7), opts)
	var buf bytes.Buffer
	heroom.WriteText(&buf, r)
	out := buf.String()
	for _, want := range []string{"Headroom", "D "Drift rate", "Count headroom", "Rate headroom"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	r := headroom.Compute(makeResults(1, 4), headroom.Options{MaxDrifted: 3})
	var buf bytes.Buffer
	if err := headroom.WriteJSON(&buf, r); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["count_headroom"]; !ok {
		t.Error("expected count_headroom field in JSON")
	}
}

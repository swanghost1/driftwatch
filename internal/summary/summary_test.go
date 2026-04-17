package summary_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/driftwatch/driftwatch/internal/summary"
)

func makeResults(drifted, clean int) []drift.Result {
	var out []drift.Result
	for i := 0; i < drifted; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < clean; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: false})
	}
	return out
}

func TestBuild_Counts(t *testing.T) {
	r := summary.Build(makeResults(3, 2))
	if r.Total != 5 {
		t.Fatalf("expected total 5, got %d", r.Total)
	}
	if r.Drifted != 3 {
		t.Fatalf("expected drifted 3, got %d", r.Drifted)
	}
	if r.Clean != 2 {
		t.Fatalf("expected clean 2, got %d", r.Clean)
	}
}

func TestBuild_DriftRate(t *testing.T) {
	r := summary.Build(makeResults(1, 3))
	if r.DriftRate != 25.0 {
		t.Fatalf("expected drift rate 25.0, got %.1f", r.DriftRate)
	}
}

func TestBuild_Empty(t *testing.T) {
	r := summary.Build(nil)
	if r.Total != 0 || r.DriftRate != 0 {
		t.Fatal("expected zero values for empty input")
	}
}

func TestWriteText_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	r := summary.Build(makeResults(2, 8))
	summary.WriteText(&buf, r)
	out := buf.String()
	for _, want := range []string{"Run at", "Total", "Drifted", "Clean", "Drift rate"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := summary.Build(makeResults(1, 1))
	if err := summary.WriteJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got summary.Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Total != 2 {
		t.Fatalf("expected total 2, got %d", got.Total)
	}
}

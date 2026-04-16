package pivot_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/pivot"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true},
		{Service: "api", Field: "replicas", Drifted: false},
		{Service: "worker", Field: "image", Drifted: true},
		{Service: "worker", Field: "env", Drifted: true},
		{Service: "cache", Field: "image", Drifted: false},
	}
}

func TestCompute_ByService_KeyCount(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisService)
	if len(tbl.Cells) != 3 {
		t.Fatalf("expected 3 service cells, got %d", len(tbl.Cells))
	}
}

func TestCompute_ByField_KeyCount(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisField)
	if len(tbl.Cells) != 3 {
		t.Fatalf("expected 3 field cells, got %d", len(tbl.Cells))
	}
}

func TestCompute_ByService_Totals(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisService)
	for _, c := range tbl.Cells {
		if c.Key == "worker" {
			if c.Total != 2 || c.Drifted != 2 || c.Clean != 0 {
				t.Errorf("worker: got total=%d drifted=%d clean=%d", c.Total, c.Drifted, c.Clean)
			}
			return
		}
	}
	t.Fatal("worker cell not found")
}

func TestCompute_ByService_DriftPct(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisService)
	for _, c := range tbl.Cells {
		if c.Key == "api" {
			if c.DriftPct != 50.0 {
				t.Errorf("api drift pct: expected 50.0, got %.1f", c.DriftPct)
			}
			return
		}
	}
	t.Fatal("api cell not found")
}

func TestCompute_SortedByDriftedDesc(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisService)
	if tbl.Cells[0].Key != "worker" {
		t.Errorf("expected worker first (most drifted), got %s", tbl.Cells[0].Key)
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	tbl := pivot.Compute(nil, pivot.AxisService)
	if len(tbl.Cells) != 0 {
		t.Errorf("expected no cells for empty input")
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisService)
	var buf bytes.Buffer
	if err := pivot.Write(&buf, tbl); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "DRIFTED") {
		t.Errorf("output missing header: %s", out)
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	tbl := pivot.Compute(makeResults(), pivot.AxisField)
	var buf bytes.Buffer
	if err := pivot.WriteJSON(&buf, tbl); err != nil {
		t.Fatal(err)
	}
	var decoded pivot.Table
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Axis != pivot.AxisField {
		t.Errorf("expected axis field, got %s", decoded.Axis)
	}
}

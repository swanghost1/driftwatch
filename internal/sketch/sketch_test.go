package sketch_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/sketch"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true},
		{Service: "api", Field: "replicas", Drifted: false},
		{Service: "worker", Field: "image", Drifted: false},
		{Service: "worker", Field: "env.LOG_LEVEL", Drifted: true},
		{Service: "cache", Field: "image", Drifted: false},
	}
}

func TestCompute_ServiceCount(t *testing.T) {
	entries := sketch.Compute(makeResults())
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestCompute_DriftedFlagSet(t *testing.T) {
	entries := sketch.Compute(makeResults())
	for _, e := range entries {
		switch e.Service {
		case "api", "worker":
			if !e.Drifted {
				t.Errorf("%s should be drifted", e.Service)
			}
		case "cache":
			if e.Drifted {
				t.Errorf("cache should not be drifted")
			}
		}
	}
}

func TestCompute_FieldsListed(t *testing.T) {
	entries := sketch.Compute(makeResults())
	for _, e := range entries {
		if e.Service == "api" {
			if len(e.Fields) != 1 || e.Fields[0] != "image" {
				t.Errorf("unexpected fields for api: %v", e.Fields)
			}
		}
	}
}

func TestCompute_Totals(t *testing.T) {
	entries := sketch.Compute(makeResults())
	for _, e := range entries {
		if e.Service == "api" && e.Total != 2 {
			t.Errorf("expected total 2 for api, got %d", e.Total)
		}
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	entries := sketch.Compute(nil)
	if len(entries) != 0 {
		t.Fatalf("expected empty, got %d", len(entries))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	sketch.Write(&buf, sketch.Compute(makeResults()))
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header in output")
	}
}

func TestWrite_ShowsDriftedYES(t *testing.T) {
	var buf bytes.Buffer
	sketch.Write(&buf, sketch.Compute(makeResults()))
	if !strings.Contains(buf.String(), "YES") {
		t.Error("expected YES for drifted service")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := sketch.WriteJSON(&buf, sketch.Compute(makeResults())); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out []sketch.Entry
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries in JSON, got %d", len(out))
	}
}

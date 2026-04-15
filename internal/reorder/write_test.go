package reorder_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/reorder"
)

func TestWriteText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	reorder.WriteText(&buf, nil)
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Fatalf("expected header in output, got: %q", buf.String())
	}
}

func TestWriteText_DriftedRow_ShowsYES(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Field: "image", Drifted: true, Expected: "nginx:1"},
	}
	var buf bytes.Buffer
	reorder.WriteText(&buf, results)
	if !strings.Contains(buf.String(), "YES") {
		t.Fatalf("expected YES in drifted row, got: %q", buf.String())
	}
}

func TestWriteText_CleanRow_ShowsNo(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Field: "replicas", Drifted: false, Expected: "3"},
	}
	var buf bytes.Buffer
	reorder.WriteText(&buf, results)
	if !strings.Contains(buf.String(), "no") {
		t.Fatalf("expected 'no' in clean row, got: %q", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	results := []drift.Result{
		{Service: "svc", Field: "image", Drifted: true, Expected: "nginx:1"},
	}
	var buf bytes.Buffer
	if err := reorder.WriteJSON(&buf, results); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestWriteJSON_EmptyResults_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := reorder.WriteJSON(&buf, nil); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "null") && !strings.Contains(buf.String(), "[]") {
		t.Fatalf("unexpected empty output: %q", buf.String())
	}
}

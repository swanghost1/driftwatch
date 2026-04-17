package peak_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/peak"
)

func makeDrifted(service, field string) drift.Result {
	return drift.Result{Service: service, Field: field, Drifted: true}
}

func makeClean(service string) drift.Result {
	return drift.Result{Service: service, Field: "image", Drifted: false}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	got := peak.Apply(nil, peak.DefaultOptions())
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestApply_OrderedByDriftCountDesc(t *testing.T) {
	runs := map[string][]drift.Result{
		"run-a": {makeDrifted("svc", "image")},
		"run-b": {makeDrifted("svc", "image"), makeDrifted("svc", "replicas")},
		"run-c": {makeClean("svc")},
	}
	got := peak.Apply(runs, peak.Options{TopN: 0})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].RunID != "run-b" {
		t.Errorf("expected run-b first, got %s", got[0].RunID)
	}
	if got[0].DriftCount != 2 {
		t.Errorf("expected drift count 2, got %d", got[0].DriftCount)
	}
}

func TestApply_TopN_LimitsResults(t *testing.T) {
	runs := map[string][]drift.Result{
		"r1": {makeDrifted("a", "image"), makeDrifted("b", "image")},
		"r2": {makeDrifted("a", "image")},
		"r3": {makeClean("a")},
	}
	got := peak.Apply(runs, peak.Options{TopN: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_NoDrift_ZeroCount(t *testing.T) {
	runs := map[string][]drift.Result{
		"run-x": {makeClean("svc")},
	}
	got := peak.Apply(runs, peak.DefaultOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].DriftCount != 0 {
		t.Errorf("expected 0 drift, got %d", got[0].DriftCount)
	}
}

func TestWrite_ContainsRunID(t *testing.T) {
	entries := []peak.Entry{{RunID: "abc-123", DriftCount: 3}}
	var buf bytes.Buffer
	peak.Write(&buf, entries)
	if !bytes.Contains(buf.Bytes(), []byte("abc-123")) {
		t.Errorf("expected run id in output, got: %s", buf.String())
	}
}

func TestWrite_EmptyEntries_ShowsMessage(t *testing.T) {
	var buf bytes.Buffer
	peak.Write(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("no peak")) {
		t.Errorf("expected 'no peak' message, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	entries := []peak.Entry{{RunID: "r1", DriftCount: 1}}
	var buf bytes.Buffer
	if err := peak.WriteJSON(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []peak.Entry
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) != 1 || out[0].RunID != "r1" {
		t.Errorf("unexpected output: %+v", out)
	}
}

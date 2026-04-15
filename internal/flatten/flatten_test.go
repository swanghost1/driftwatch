package flatten_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/flatten"
)

func makeResults(drifted bool) []flatten.Result {
	return []flatten.Result{
		{Service: "api", Field: "image", Wanted: "nginx:1.25", Actual: "nginx:1.24", Drifted: drifted},
		{Service: "api", Field: "replicas", Wanted: "3", Actual: "2", Drifted: drifted},
		{Service: "worker", Field: "image", Wanted: "alpine:3.18", Actual: "alpine:3.18", Drifted: false},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	results := makeResults(true)
	rows := flatten.Apply(results)
	if len(rows) != len(results) {
		t.Fatalf("expected %d rows, got %d", len(results), len(rows))
	}
}

func TestApply_SortedByServiceThenField(t *testing.T) {
	results := []flatten.Result{
		{Service: "worker", Field: "replicas"},
		{Service: "api", Field: "image"},
		{Service: "api", Field: "env"},
	}
	rows := flatten.Apply(results)
	if rows[0].Service != "api" || rows[0].Field != "env" {
		t.Errorf("expected api/env first, got %s/%s", rows[0].Service, rows[0].Field)
	}
	if rows[1].Service != "api" || rows[1].Field != "image" {
		t.Errorf("expected api/image second, got %s/%s", rows[1].Service, rows[1].Field)
	}
	if rows[2].Service != "worker" {
		t.Errorf("expected worker last, got %s", rows[2].Service)
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	rows := flatten.Apply(nil)
	if len(rows) != 0 {
		t.Fatalf("expected empty slice, got %d rows", len(rows))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	rows := flatten.Apply(makeResults(true))
	if err := flatten.Write(&buf, rows); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"SERVICE", "FIELD", "WANTED", "ACTUAL", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWrite_DriftedRow_ShowsDRIFT(t *testing.T) {
	var buf bytes.Buffer
	rows := flatten.Apply([]flatten.Result{
		{Service: "svc", Field: "image", Wanted: "a", Actual: "b", Drifted: true},
	})
	_ = flatten.Write(&buf, rows)
	if !strings.Contains(buf.String(), "DRIFT") {
		t.Error("expected DRIFT status in output")
	}
}

func TestWrite_CleanRow_ShowsOK(t *testing.T) {
	var buf bytes.Buffer
	rows := flatten.Apply([]flatten.Result{
		{Service: "svc", Field: "image", Wanted: "a", Actual: "a", Drifted: false},
	})
	_ = flatten.Write(&buf, rows)
	if !strings.Contains(buf.String(), "OK") {
		t.Error("expected OK status in output")
	}
}

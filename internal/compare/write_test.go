package compare_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/compare"
	"github.com/example/driftwatch/internal/drift"
)

func TestWriteText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	if err := compare.WriteText(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' message, got: %s", buf.String())
	}
}

func TestWriteText_ContainsHeaders(t *testing.T) {
	changes := []compare.Change{
		{
			Result: drift.Result{Service: "svc-a", Field: "image", Expected: "v1", Actual: "v2", Drifted: true},
			Kind:   compare.ChangeIntroduced,
		},
	}
	var buf bytes.Buffer
	if err := compare.WriteText(&buf, changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, header := range []string{"SERVICE", "FIELD", "EXPECTED", "ACTUAL", "CHANGE"} {
		if !strings.Contains(out, header) {
			t.Errorf("missing header %q in output: %s", header, out)
		}
	}
}

func TestWriteText_ContainsChangeKind(t *testing.T) {
	changes := []compare.Change{
		{
			Result: drift.Result{Service: "svc-b", Field: "replicas", Expected: "3", Actual: "1", Drifted: true},
			Kind:   compare.ChangeResolved,
		},
	}
	var buf bytes.Buffer
	_ = compare.WriteText(&buf, changes)
	if !strings.Contains(buf.String(), string(compare.ChangeResolved)) {
		t.Errorf("expected 'resolved' in output, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	changes := []compare.Change{
		{
			Result: drift.Result{Service: "svc-c", Field: "env.PORT", Expected: "8080", Actual: "9090", Drifted: true},
			Kind:   compare.ChangeIntroduced,
		},
	}
	var buf bytes.Buffer
	if err := compare.WriteJSON(&buf, changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["service"] != "svc-c" {
		t.Errorf("expected service svc-c, got %s", rows[0]["service"])
	}
	if rows[0]["change"] != string(compare.ChangeIntroduced) {
		t.Errorf("expected change introduced, got %s", rows[0]["change"])
	}
}

func TestWriteJSON_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	if err := compare.WriteJSON(&buf, []compare.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array, got: %s", buf.String())
	}
}

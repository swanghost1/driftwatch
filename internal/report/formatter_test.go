package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{
			ServiceName: "api",
			Diffs:       []drift.Diff{},
		},
		{
			ServiceName: "worker",
			Diffs: []drift.Diff{
				{Field: "image", Expected: "app:1.2", Actual: "app:1.1"},
				{Field: "replicas", Expected: "3", Actual: "2"},
			},
		},
	}
}

func TestNewSummary_Counts(t *testing.T) {
	s := report.NewSummary(makeResults())
	if s.TotalServices != 2 {
		t.Errorf("expected TotalServices=2, got %d", s.TotalServices)
	}
	if s.DriftedServices != 1 {
		t.Errorf("expected DriftedServices=1, got %d", s.DriftedServices)
	}
}

func TestWrite_TextFormat_ContainsOKAndDrift(t *testing.T) {
	var buf bytes.Buffer
	s := report.NewSummary(makeResults())
	if err := report.Write(&buf, s, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]    api") {
		t.Errorf("expected OK line for api, got:\n%s", out)
	}
	if !strings.Contains(out, "[DRIFT] worker") {
		t.Errorf("expected DRIFT line for worker, got:\n%s", out)
	}
	if !strings.Contains(out, "image") {
		t.Errorf("expected image diff in output, got:\n%s", out)
	}
}

func TestWrite_TextFormat_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	s := report.NewSummary([]drift.Result{{ServiceName: "db", Diffs: []drift.Diff{}}})
	if err := report.Write(&buf, s, report.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "[DRIFT]") {
		t.Errorf("expected no DRIFT lines when no drift present")
	}
}

func TestWrite_JSONFormat_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	s := report.NewSummary(makeResults())
	if err := report.Write(&buf, s, report.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"total_services", "drifted_services", "results", "worker"} {
		if !strings.Contains(out, want) {
			t.Errorf("JSON output missing %q:\n%s", want, out)
		}
	}
}

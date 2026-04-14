package aggregate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/aggregate"
	"github.com/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true},
		{Service: "worker", Field: "image", Drifted: false},
		{Service: "api", Field: "replicas", Drifted: false},
		{Service: "worker", Field: "replicas", Drifted: true},
		{Service: "gateway", Field: "image", Drifted: true},
	}
}

func TestByField_FieldCount(t *testing.T) {
	summaries := aggregate.ByField(makeResults())
	if len(summaries) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(summaries))
	}
}

func TestByField_ImageTotals(t *testing.T) {
	summaries := aggregate.ByField(makeResults())
	var img aggregate.FieldSummary
	for _, s := range summaries {
		if s.Field == "image" {
			img = s
		}
	}
	if img.Total != 3 {
		t.Errorf("expected total 3, got %d", img.Total)
	}
	if img.Drifted != 2 {
		t.Errorf("expected drifted 2, got %d", img.Drifted)
	}
}

func TestByField_DriftRate(t *testing.T) {
	summaries := aggregate.ByField(makeResults())
	for _, s := range summaries {
		if s.Field == "image" {
			rate := s.DriftRate()
			if rate < 0.66 || rate > 0.67 {
				t.Errorf("unexpected drift rate %.4f for image", rate)
			}
		}
	}
}

func TestByField_SortedAlphabetically(t *testing.T) {
	summaries := aggregate.ByField(makeResults())
	if summaries[0].Field != "image" || summaries[1].Field != "replicas" {
		t.Errorf("expected sorted order image, replicas; got %s, %s",
			summaries[0].Field, summaries[1].Field)
	}
}

func TestByField_EmptyInput(t *testing.T) {
	summaries := aggregate.ByField(nil)
	if len(summaries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(summaries))
	}
}

func TestByField_ServicesListedPerField(t *testing.T) {
	summaries := aggregate.ByField(makeResults())
	for _, s := range summaries {
		if s.Field == "image" && len(s.Services) != 3 {
			t.Errorf("expected 3 services for image, got %d", len(s.Services))
		}
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	summaries := aggregate.ByField(makeResults())
	if err := aggregate.Write(&buf, summaries); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"FIELD", "TOTAL", "DRIFTED", "RATE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
}

func TestWrite_ContainsFieldRow(t *testing.T) {
	var buf bytes.Buffer
	summaries := aggregate.ByField(makeResults())
	_ = aggregate.Write(&buf, summaries)
	if !strings.Contains(buf.String(), "image") {
		t.Error("output missing 'image' row")
	}
}

package coalesce_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/coalesce"
	"github.com/driftwatch/internal/drift"
)

func makeResult(service, field string, drifted bool, t time.Time) drift.Result {
	return drift.Result{
		Service:    service,
		Field:      field,
		Drifted:    drifted,
		DetectedAt: t,
	}
}

var (
	t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = t0.Add(time.Hour)
)

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	got := coalesce.Apply(nil, coalesce.DefaultOptions())
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestApply_NoDuplicates_ReturnsAll(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", true, t0),
		makeResult("svc-b", "image", false, t0),
	}
	got := coalesce.Apply(input, coalesce.DefaultOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_DuplicateClean_PreferDrifted_KeepsDrifted(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", false, t0),
		makeResult("svc-a", "image", true, t0),
	}
	opts := coalesce.DefaultOptions()
	got := coalesce.Apply(input, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if !got[0].Drifted {
		t.Error("expected drifted result to be kept")
	}
}

func TestApply_DuplicateDrifted_PreferNewest_KeepsNewer(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "replicas", true, t0),
		makeResult("svc-a", "replicas", true, t1),
	}
	opts := coalesce.DefaultOptions()
	got := coalesce.Apply(input, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if !got[0].DetectedAt.Equal(t1) {
		t.Errorf("expected newer timestamp %v, got %v", t1, got[0].DetectedAt)
	}
}

func TestApply_SameServiceDifferentFields_BothKept(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", true, t0),
		makeResult("svc-a", "replicas", true, t0),
	}
	got := coalesce.Apply(input, coalesce.DefaultOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApply_SortedByServiceThenField(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-b", "image", false, t0),
		makeResult("svc-a", "replicas", true, t0),
		makeResult("svc-a", "image", true, t0),
	}
	got := coalesce.Apply(input, coalesce.DefaultOptions())
	if got[0].Service != "svc-a" || got[0].Field != "image" {
		t.Errorf("unexpected first result: %+v", got[0])
	}
	if got[1].Field != "replicas" {
		t.Errorf("unexpected second result: %+v", got[1])
	}
}

func TestWrite_ContainsSummaryLine(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", true, t0),
		makeResult("svc-b", "image", false, t0),
	}
	var buf bytes.Buffer
	coalesce.Write(&buf, input)
	out := buf.String()
	if !strings.Contains(out, "coalesced: 2 results") {
		t.Errorf("missing summary line, got: %s", out)
	}
	if !strings.Contains(out, "DRIFT") {
		t.Error("expected DRIFT label in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK label in output")
	}
}

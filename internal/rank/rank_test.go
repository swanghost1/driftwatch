package rank_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/rank"
)

func makeResult(service, field, severity string, drifted bool) drift.Result {
	return drift.Result{
		Service:  service,
		Field:    field,
		Severity: severity,
		Drifted:  drifted,
	}
}

func TestApply_DriftedFirst(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", "low", false),
		makeResult("svc-b", "image", "critical", true),
		makeResult("svc-c", "replicas", "high", false),
	}
	opts := rank.Options{Criteria: []rank.Criterion{rank.ByDriftOnly}}
	out := rank.Apply(input, opts)

	if !out[0].Drifted {
		t.Errorf("expected first result to be drifted, got service=%s", out[0].Service)
	}
}

func TestApply_BySeverity_CriticalFirst(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", "low", true),
		makeResult("svc-b", "image", "critical", true),
		makeResult("svc-c", "image", "high", true),
	}
	opts := rank.Options{Criteria: []rank.Criterion{rank.BySeverity}}
	out := rank.Apply(input, opts)

	if out[0].Severity != "critical" {
		t.Errorf("expected critical first, got %s", out[0].Severity)
	}
	if out[1].Severity != "high" {
		t.Errorf("expected high second, got %s", out[1].Severity)
	}
}

func TestApply_ByService_Alphabetical(t *testing.T) {
	input := []drift.Result{
		makeResult("zebra", "image", "", false),
		makeResult("alpha", "image", "", false),
		makeResult("mango", "image", "", false),
	}
	opts := rank.Options{Criteria: []rank.Criterion{rank.ByService}}
	out := rank.Apply(input, opts)

	if out[0].Service != "alpha" {
		t.Errorf("expected alpha first, got %s", out[0].Service)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-z", "image", "low", false),
		makeResult("svc-a", "image", "critical", true),
	}
	origFirst := input[0].Service
	rank.Apply(input, rank.DefaultOptions())
	if input[0].Service != origFirst {
		t.Error("Apply must not mutate the original slice")
	}
}

func TestApply_DefaultOptions_DriftedAndSeverityOrdered(t *testing.T) {
	input := []drift.Result{
		makeResult("svc-a", "image", "low", false),
		makeResult("svc-b", "image", "high", true),
		makeResult("svc-c", "image", "critical", true),
	}
	out := rank.Apply(input, rank.DefaultOptions())

	if !out[0].Drifted || !out[1].Drifted {
		t.Error("expected drifted results to appear first")
	}
	if out[0].Severity != "critical" {
		t.Errorf("expected critical before high, got %s", out[0].Severity)
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	results := []drift.Result{
		makeResult("my-service", "image", "critical", true),
	}
	rank.Write(&buf, results)
	out := buf.String()

	for _, hdr := range []string{"SERVICE", "FIELD", "SEVERITY", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestWrite_ShowsDriftStatus(t *testing.T) {
	var buf bytes.Buffer
	results := []drift.Result{
		makeResult("my-service", "image", "critical", true),
		makeResult("clean-svc", "replicas", "low", false),
	}
	rank.Write(&buf, results)
	out := buf.String()

	if !strings.Contains(out, "DRIFT") {
		t.Error("expected DRIFT label in output")
	}
	if !strings.Contains(out, "ok") {
		t.Error("expected ok label in output")
	}
}

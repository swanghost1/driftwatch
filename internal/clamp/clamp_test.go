package clamp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/clamp"
	"github.com/driftwatch/internal/drift"
)

func makeReplicaResult(service, live, want string) drift.Result {
	return drift.Result{
		Service:   service,
		Field:     "replicas",
		LiveValue: live,
		WantValue: want,
		Drifted:   live != want,
	}
}

func makeImageResult(service string) drift.Result {
	return drift.Result{
		Service:   service,
		Field:     "image",
		LiveValue: "nginx:1.24",
		WantValue: "nginx:1.25",
		Drifted:   true,
	}
}

func TestApply_NoViolations_WithinBounds(t *testing.T) {
	results := []drift.Result{makeReplicaResult("svc-a", "3", "3")}
	opts := clamp.Options{MinReplicas: 1, MaxReplicas: 5}
	got := clamp.Apply(results, opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(got))
	}
}

func TestApply_BelowMinimum_ReturnsViolation(t *testing.T) {
	results := []drift.Result{makeReplicaResult("svc-a", "0", "2")}
	opts := clamp.Options{MinReplicas: 1, MaxReplicas: 10}
	got := clamp.Apply(results, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Service != "svc-a" {
		t.Errorf("unexpected service %q", got[0].Service)
	}
	if !strings.Contains(got[0].Description, "below minimum") {
		t.Errorf("description missing 'below minimum': %s", got[0].Description)
	}
}

func TestApply_AboveMaximum_ReturnsViolation(t *testing.T) {
	results := []drift.Result{makeReplicaResult("svc-b", "20", "4")}
	opts := clamp.Options{MinReplicas: 1, MaxReplicas: 10}
	got := clamp.Apply(results, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if !strings.Contains(got[0].Description, "above maximum") {
		t.Errorf("description missing 'above maximum': %s", got[0].Description)
	}
}

func TestApply_NonReplicaField_Ignored(t *testing.T) {
	results := []drift.Result{makeImageResult("svc-c")}
	opts := clamp.Options{MinReplicas: 1, MaxReplicas: 5}
	got := clamp.Apply(results, opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 violations for image field, got %d", len(got))
	}
}

func TestApply_ZeroBounds_NoBoundChecked(t *testing.T) {
	results := []drift.Result{makeReplicaResult("svc-d", "100", "2")}
	opts := clamp.Options{MinReplicas: 0, MaxReplicas: 0}
	got := clamp.Apply(results, opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 violations when bounds are zero, got %d", len(got))
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	violations := []clamp.Violation{
		{Service: "svc-a", LiveValue: 0, Declared: 2, Min: 1, Max: 5, Description: "live replicas 0 below minimum 1"},
	}
	var buf bytes.Buffer
	clamp.Write(&buf, violations)
	out := buf.String()
	for _, hdr := range []string{"SERVICE", "LIVE", "DECLARED", "MIN", "MAX", "DESCRIPTION"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
}

func TestWrite_ContainsServiceName(t *testing.T) {
	violations := []clamp.Violation{
		{Service: "my-service", LiveValue: 12, Declared: 3, Min: 1, Max: 10, Description: "above maximum"},
	}
	var buf bytes.Buffer
	clamp.Write(&buf, violations)
	if !strings.Contains(buf.String(), "my-service") {
		t.Error("output does not contain service name")
	}
}

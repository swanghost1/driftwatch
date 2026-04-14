package enrich_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/enrich"
)

func makeResults(drifted bool) []enrich.Result {
	return []enrich.Result{
		{
			Service:  "api",
			Field:    "image",
			Expected: "nginx:1.25",
			Actual:   "nginx:1.24",
			Drifted:  drifted,
		},
	}
}

func TestApply_StampsDetectedAt(t *testing.T) {
	results := makeResults(true)
	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	enrich.Apply(results, enrich.Options{DetectedAt: fixed})
	if !results[0].DetectedAt.Equal(fixed) {
		t.Errorf("expected DetectedAt %v, got %v", fixed, results[0].DetectedAt)
	}
}

func TestApply_DoesNotOverwriteExistingTime(t *testing.T) {
	existing := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	results := []enrich.Result{
		{Service: "svc", Field: "replicas", Drifted: true, DetectedAt: existing},
	}
	newer := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	enrich.Apply(results, enrich.Options{DetectedAt: newer})
	if !results[0].DetectedAt.Equal(existing) {
		t.Errorf("DetectedAt should not be overwritten; got %v", results[0].DetectedAt)
	}
}

func TestApply_AddsDescription_WhenDrifted(t *testing.T) {
	results := makeResults(true)
	enrich.Apply(results, enrich.Options{})
	if results[0].Description == "" {
		t.Error("expected a description to be set for drifted result")
	}
	if !strings.Contains(results[0].Description, "nginx:1.25") {
		t.Errorf("description should mention expected value; got %q", results[0].Description)
	}
	if !strings.Contains(results[0].Description, "nginx:1.24") {
		t.Errorf("description should mention actual value; got %q", results[0].Description)
	}
}

func TestApply_NoDescription_WhenNotDrifted(t *testing.T) {
	results := makeResults(false)
	enrich.Apply(results, enrich.Options{})
	if results[0].Description != "" {
		t.Errorf("expected no description for clean result, got %q", results[0].Description)
	}
}

func TestApply_CustomDescriber_IsUsed(t *testing.T) {
	results := makeResults(true)
	called := false
	custom := func(field, expected, actual string) string {
		called = true
		return "custom: " + field
	}
	enrich.Apply(results, enrich.Options{Describe: custom})
	if !called {
		t.Error("expected custom describer to be called")
	}
	if !strings.HasPrefix(results[0].Description, "custom:") {
		t.Errorf("unexpected description: %q", results[0].Description)
	}
}

func TestApply_DefaultsDetectedAt_WhenZero(t *testing.T) {
	before := time.Now().UTC()
	results := makeResults(true)
	enrich.Apply(results, enrich.Options{})
	after := time.Now().UTC()
	if results[0].DetectedAt.Before(before) || results[0].DetectedAt.After(after) {
		t.Errorf("DetectedAt %v not within expected range", results[0].DetectedAt)
	}
}

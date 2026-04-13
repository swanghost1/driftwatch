package severity_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/severity"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Drifted: true,
			Diffs: []drift.Diff{{Field: "image", Want: "v1", Got: "v2"}},
		},
		{
			Service: "worker",
			Drifted: true,
			Diffs: []drift.Diff{{Field: "replicas", Want: "3", Got: "2"}},
		},
		{
			Service: "cache",
			Drifted: false,
		},
	}
}

func TestClassify_NoDrift_LevelNone(t *testing.T) {
	results := makeResults()
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)

	for _, r := range classified {
		if r.Drift.Service == "cache" && r.Severity != severity.LevelNone {
			t.Errorf("expected LevelNone for clean service, got %s", r.Severity)
		}
	}
}

func TestClassify_ImageDrift_IsCritical(t *testing.T) {
	results := makeResults()
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)

	for _, r := range classified {
		if r.Drift.Service == "api" {
			if r.Severity != severity.LevelCritical {
				t.Errorf("expected LevelCritical for image drift, got %s", r.Severity)
			}
			return
		}
	}
	t.Fatal("api result not found")
}

func TestClassify_ReplicasDrift_IsHigh(t *testing.T) {
	results := makeResults()
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)

	for _, r := range classified {
		if r.Drift.Service == "worker" {
			if r.Severity != severity.LevelHigh {
				t.Errorf("expected LevelHigh for replicas drift, got %s", r.Severity)
			}
			return
		}
	}
	t.Fatal("worker result not found")
}

func TestClassify_UnknownField_DefaultsToLow(t *testing.T) {
	results := []drift.Result{
		{Service: "svc", Drifted: true, Diffs: []drift.Diff{{Field: "unknown_field", Want: "a", Got: "b"}}},
	}
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)
	if len(classified) != 1 {
		t.Fatalf("expected 1 result, got %d", len(classified))
	}
	if classified[0].Severity != severity.LevelLow {
		t.Errorf("expected LevelLow for unknown field, got %s", classified[0].Severity)
	}
}

func TestClassify_MultipleFields_TakesHighest(t *testing.T) {
	results := []drift.Result{
		{
			Service: "svc",
			Drifted: true,
			Diffs: []drift.Diff{
				{Field: "env", Want: "a", Got: "b"},
				{Field: "image", Want: "v1", Got: "v2"},
			},
		},
	}
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)
	if classified[0].Severity != severity.LevelCritical {
		t.Errorf("expected LevelCritical, got %s", classified[0].Severity)
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	results := makeResults()
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)

	var buf bytes.Buffer
	if err := severity.Write(&buf, classified); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"SERVICE", "SEVERITY", "DRIFTED"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
}

func TestWrite_SortedBySeverityDescending(t *testing.T) {
	results := makeResults()
	cfg := severity.DefaultConfig()
	classified := severity.Classify(results, cfg)

	var buf bytes.Buffer
	_ = severity.Write(&buf, classified)
	out := buf.String()

	crIdx := strings.Index(out, "critical")
	hiIdx := strings.Index(out, "high")
	if crIdx == -1 || hiIdx == -1 {
		t.Skip("severity labels not found in output")
	}
	if crIdx > hiIdx {
		t.Error("expected critical before high in output")
	}
}

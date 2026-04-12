package remediate_test

import (
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/remediate"
)

func makeResults(drifted bool, diffs []drift.Diff) []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Drifted: drifted,
			Diffs:   diffs,
		},
	}
}

func TestSuggest_NoDrift_ReturnsEmpty(t *testing.T) {
	results := makeResults(false, nil)
	got := remediate.Suggest(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(got))
	}
}

func TestSuggest_ImageDrift_ContainsKubectlSetImage(t *testing.T) {
	results := makeResults(true, []drift.Diff{
		{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.21"},
	})
	got := remediate.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Action, "kubectl set image") {
		t.Errorf("expected kubectl set image, got: %s", got[0].Action)
	}
	if !strings.Contains(got[0].Action, "nginx:1.25") {
		t.Errorf("expected expected value in action, got: %s", got[0].Action)
	}
}

func TestSuggest_ReplicasDrift_ContainsKubectlScale(t *testing.T) {
	results := makeResults(true, []drift.Diff{
		{Field: "replicas", Expected: 3, Actual: 1},
	})
	got := remediate.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Action, "kubectl scale") {
		t.Errorf("expected kubectl scale, got: %s", got[0].Action)
	}
}

func TestSuggest_EnvDrift_ContainsKubectlSetEnv(t *testing.T) {
	results := makeResults(true, []drift.Diff{
		{Field: "env:LOG_LEVEL", Expected: "info", Actual: "debug"},
	})
	got := remediate.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Action, "kubectl set env") {
		t.Errorf("expected kubectl set env, got: %s", got[0].Action)
	}
}

func TestSuggest_UnknownField_ContainsManualHint(t *testing.T) {
	results := makeResults(true, []drift.Diff{
		{Field: "cpu_limit", Expected: "500m", Actual: "250m"},
	})
	got := remediate.Suggest(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Action, "manually reconcile") {
		t.Errorf("expected manual hint, got: %s", got[0].Action)
	}
}

func TestSuggest_MultipleDiffs_ReturnsSuggestionPerDiff(t *testing.T) {
	results := makeResults(true, []drift.Diff{
		{Field: "image", Expected: "app:2.0", Actual: "app:1.0"},
		{Field: "replicas", Expected: 2, Actual: 1},
	})
	got := remediate.Suggest(results)
	if len(got) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(got))
	}
}

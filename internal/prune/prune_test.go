package prune_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"driftwatch/internal/prune"
)

func makeResults() []prune.Result {
	now := time.Now()
	return []prune.Result{
		{Service: "alpha", Field: "image", Drifted: true, DetectedAt: now},
		{Service: "beta", Field: "replicas", Drifted: false, DetectedAt: now},
		{Service: "gamma", Field: "image", Drifted: true, DetectedAt: now.Add(-48 * time.Hour)},
		{Service: "delta", Field: "env", Drifted: false, DetectedAt: now.Add(-72 * time.Hour)},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := makeResults()
	out := prune.Apply(results, prune.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_OlderThan_RemovesOldEntries(t *testing.T) {
	results := makeResults()
	out := prune.Apply(results, prune.Options{OlderThan: 24 * time.Hour})
	// gamma (48h) and delta (72h) should be pruned
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Service == "gamma" || r.Service == "delta" {
			t.Errorf("expected %q to be pruned", r.Service)
		}
	}
}

func TestApply_OnlyClean_RemovesNonDrifted(t *testing.T) {
	results := makeResults()
	out := prune.Apply(results, prune.Options{OnlyClean: true})
	for _, r := range out {
		if !r.Drifted {
			t.Errorf("expected only drifted results, got clean result for %q", r.Service)
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(out))
	}
}

func TestApply_ServiceFilter_RemovesMatchingService(t *testing.T) {
	results := makeResults()
	out := prune.Apply(results, prune.Options{Services: []string{"alpha", "beta"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Service == "alpha" || r.Service == "beta" {
			t.Errorf("expected %q to be pruned", r.Service)
		}
	}
}

func TestApply_ServiceFilter_CaseInsensitive(t *testing.T) {
	results := makeResults()
	out := prune.Apply(results, prune.Options{Services: []string{"ALPHA"}})
	for _, r := range out {
		if r.Service == "alpha" {
			t.Error("expected alpha to be pruned (case-insensitive)")
		}
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := prune.Apply([]prune.Result{}, prune.Options{OnlyClean: true})
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d", len(out))
	}
}

func TestWrite_ShowsCorrectCounts(t *testing.T) {
	before := makeResults()
	after := before[:2]
	var buf bytes.Buffer
	prune.Write(&buf, before, after)
	got := buf.String()
	if !strings.Contains(got, "4") || !strings.Contains(got, "2 removed") {
		t.Errorf("unexpected output: %q", got)
	}
}

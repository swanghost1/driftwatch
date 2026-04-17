package threshold_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/threshold"
)

func makeResults(drifted, clean int) []drift.Result {
	var out []drift.Result
	for i := 0; i < drifted; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < clean; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: false})
	}
	return out
}

func TestEvaluate_NoDrift_NotBreached(t *testing.T) {
	r := threshold.Evaluate(makeResults(0, 5), threshold.Options{MinDriftCount: 1})
	if r.Breached {
		t.Fatal("expected no breach")
	}
}

func TestEvaluate_MeetsCountThreshold_Breached(t *testing.T) {
	r := threshold.Evaluate(makeResults(3, 2), threshold.Options{MinDriftCount: 3})
	if !r.Breached {
		t.Fatal("expected breach")
	}
	if r.Drifted != 3 {
		t.Fatalf("want 3 drifted, got %d", r.Drifted)
	}
}

func TestEvaluate_BelowCountThreshold_NotBreached(t *testing.T) {
	r := threshold.Evaluate(makeResults(2, 3), threshold.Options{MinDriftCount: 3})
	if r.Breached {
		t.Fatal("expected no breach")
	}
}

func TestEvaluate_MeetsRateThreshold_Breached(t *testing.T) {
	r := threshold.Evaluate(makeResults(4, 1), threshold.Options{MinDriftRate: 0.8})
	if !r.Breached {
		t.Fatalf("expected breach, rate=%.2f", r.DriftRate)
	}
}

func TestEvaluate_BelowRateThreshold_NotBreached(t *testing.T) {
	r := threshold.Evaluate(makeResults(1, 4), threshold.Options{MinDriftRate: 0.5})
	if r.Breached {
		t.Fatal("expected no breach")
	}
}

func TestEvaluate_EmptyResults_NotBreached(t *testing.T) {
	r := threshold.Evaluate(nil, threshold.Options{MinDriftCount: 1, MinDriftRate: 0.5})
	if r.Breached {
		t.Fatal("expected no breach on empty input")
	}
	if r.DriftRate != 0 {
		t.Fatal("expected zero drift rate")
	}
}

func TestWrite_BreachedOutput(t *testing.T) {
	r := threshold.Result{Total: 5, Drifted: 3, DriftRate: 0.6, Breached: true, Reason: "count"}
	var sb strings.Builder
	threshold.Write(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "BREACHED") {
		t.Error("expected BREACHED in output")
	}
	if !strings.Contains(out, "count") {
		t.Error("expected reason in output")
	}
}

func TestWrite_OKOutput(t *testing.T) {
	r := threshold.Result{Total: 5, Drifted: 0, DriftRate: 0}
	var sb strings.Builder
	threshold.Write(&sb, r)
	if !strings.Contains(sb.String(), "OK") {
		t.Error("expected OK in output")
	}
}

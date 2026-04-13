package variance_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/variance"
)

func makeDrifted(service string) drift.Result {
	return drift.Result{Service: service, Drifted: true, Field: "image"}
}

func makeClean(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestCompute_EmptyHistory_ReturnsNil(t *testing.T) {
	got := variance.Compute(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCompute_SingleRun_StdDevZero(t *testing.T) {
	history := [][]drift.Result{
		{makeDrifted("svc-a"), makeDrifted("svc-a")},
	}
	got := variance.Compute(history)
	if len(got) != 1 {
		t.Fatalf("expected 1 variance entry, got %d", len(got))
	}
	if got[0].StdDev != 0 {
		t.Errorf("expected stddev 0 for single sample, got %f", got[0].StdDev)
	}
	if got[0].Mean != 2 {
		t.Errorf("expected mean 2, got %f", got[0].Mean)
	}
}

func TestCompute_MultipleRuns_MeanCorrect(t *testing.T) {
	history := [][]drift.Result{
		{makeDrifted("svc-a")},
		{makeDrifted("svc-a"), makeDrifted("svc-a")},
		{makeDrifted("svc-a"), makeDrifted("svc-a"), makeDrifted("svc-a")},
	}
	got := variance.Compute(history)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if got[0].Mean != 2.0 {
		t.Errorf("expected mean 2.0, got %f", got[0].Mean)
	}
	if got[0].Samples != 3 {
		t.Errorf("expected 3 samples, got %d", got[0].Samples)
	}
}

func TestCompute_AnomalyDetected(t *testing.T) {
	// runs 1-3 have 0 drifts, run 4 has 10 — should be anomalous
	history := [][]drift.Result{
		{makeClean("svc-a")},
		{makeClean("svc-a")},
		{makeClean("svc-a")},
	}
	// add a fourth run with many drifts
	var spikeRun []drift.Result
	for i := 0; i < 10; i++ {
		spikeRun = append(spikeRun, makeDrifted("svc-a"))
	}
	history = append(history, spikeRun)

	got := variance.Compute(history)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if !got[0].Anomalous {
		t.Errorf("expected anomalous=true for spike run")
	}
}

func TestCompute_NoAnomaly_WhenBelowThreshold(t *testing.T) {
	history := [][]drift.Result{
		{makeDrifted("svc-b")},
		{makeDrifted("svc-b")},
		{makeDrifted("svc-b")},
		{makeDrifted("svc-b")},
	}
	got := variance.Compute(history)
	if got[0].Anomalous {
		t.Errorf("expected anomalous=false for stable drift")
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	vars := []variance.ServiceVariance{
		{Service: "svc-a", Mean: 1.5, StdDev: 0.5, Samples: 4, Anomalous: false},
	}
	var buf bytes.Buffer
	variance.Write(&buf, vars)
	out := buf.String()
	if !strings.Contains(out, "SERVICE") {
		t.Errorf("expected SERVICE header")
	}
	if !strings.Contains(out, "STDDEV") {
		t.Errorf("expected STDDEV header")
	}
}

func TestWrite_AnomalousFlag(t *testing.T) {
	vars := []variance.ServiceVariance{
		{Service: "svc-x", Mean: 0.1, StdDev: 0.1, Samples: 5, Anomalous: true},
	}
	var buf bytes.Buffer
	variance.Write(&buf, vars)
	if !strings.Contains(buf.String(), "YES") {
		t.Errorf("expected YES for anomalous service")
	}
}

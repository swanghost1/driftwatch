package spike

import (
	"bytes"
	"strings"
	"testing"
)

func TestDetect_NoHistory_NeverSpike(t *testing.T) {
	r := Detect(100, nil, DefaultOptions())
	if r.IsSpike {
		t.Fatal("expected no spike with empty history")
	}
}

func TestDetect_BelowThreshold_NoSpike(t *testing.T) {
	history := []int{4, 4, 4, 4, 4}
	r := Detect(5, history, DefaultOptions()) // mean=4, threshold=8
	if r.IsSpike {
		t.Fatalf("expected no spike, current=%d threshold=%.0f", r.Current, r.Threshold)
	}
}

func TestDetect_MeetsThreshold_IsSpike(t *testing.T) {
	history := []int{4, 4, 4, 4, 4}
	r := Detect(8, history, DefaultOptions()) // mean=4, threshold=8
	if !r.IsSpike {
		t.Fatalf("expected spike, current=%d threshold=%.0f", r.Current, r.Threshold)
	}
}

func TestDetect_ExceedsThreshold_IsSpike(t *testing.T) {
	history := []int{2, 2, 2}
	r := Detect(10, history, DefaultOptions()) // mean=2, threshold=4
	if !r.IsSpike {
		t.Fatal("expected spike")
	}
}

func TestDetect_WindowTruncated(t *testing.T) {
	// Only last 5 values should count; older high values ignored.
	history := []int{100, 100, 1, 1, 1, 1, 1}
	r := Detect(3, history, DefaultOptions()) // window=[1,1,1,1,1] mean=1 threshold=2
	if r.IsSpike {
		t.Fatalf("expected no spike after window truncation, threshold=%.0f", r.Threshold)
	}
}

func TestDetect_MeanCalculation(t *testing.T) {
	history := []int{2, 4, 6}
	r := Detect(1, history, DefaultOptions())
	if r.Mean != 4.0 {
		t.Fatalf("expected mean 4.0, got %.2f", r.Mean)
	}
}

func TestDetect_CustomMultiplier(t *testing.T) {
	opts := Options{WindowSize: 3, ThresholdMult: 3.0}
	history := []int{2, 2, 2}
	r := Detect(6, history, opts) // mean=2 threshold=6
	if !r.IsSpike {
		t.Fatal("expected spike with multiplier 3")
	}
}

func TestWrite_SpikeOutput(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, Result{Current: 10, Mean: 4.0, Threshold: 8, IsSpike: true})
	out := buf.String()
	if !strings.Contains(out, "YES") {
		t.Errorf("expected YES in output, got: %s", out)
	}
	if !strings.Contains(out, "10") {
		t.Errorf("expected current count in output")
	}
}

func TestWrite_NoSpike(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, Result{Current: 2, Mean: 4.0, Threshold: 8, IsSpike: false})
	out := buf.String()
	if strings.Contains(out, "YES") {
		t.Errorf("did not expect YES in output")
	}
}

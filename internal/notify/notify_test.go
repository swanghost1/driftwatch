package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/notify"
)

func makeResults(driftedNames ...string) []drift.Result {
	results := []drift.Result{
		{Service: "clean-svc", Drifted: false},
	}
	for _, name := range driftedNames {
		results = append(results, drift.Result{
			Service: name,
			Drifted: true,
			Diffs:   []drift.Diff{{Field: "image", Expected: "v1", Actual: "v2"}},
		})
	}
	return results
}

func TestNotify_NoDrift_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, notify.Options{MinDriftCount: 0})

	sent, err := n.Notify(makeResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent {
		t.Error("expected no notification when no drift present")
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestNotify_WithDrift_SendsNotification(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, notify.Options{MinDriftCount: 0})

	sent, err := n.Notify(makeResults("api-svc", "worker-svc"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sent {
		t.Error("expected notification to be sent")
	}
	out := buf.String()
	if !strings.Contains(out, "DRIFT DETECTED") {
		t.Errorf("expected DRIFT DETECTED in output, got %q", out)
	}
	if !strings.Contains(out, "api-svc") {
		t.Errorf("expected api-svc in output, got %q", out)
	}
	if !strings.Contains(out, "worker-svc") {
		t.Errorf("expected worker-svc in output, got %q", out)
	}
}

func TestNotify_BelowThreshold_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, notify.Options{MinDriftCount: 3})

	sent, err := n.Notify(makeResults("api-svc"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent {
		t.Error("expected no notification below threshold")
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output below threshold, got %q", buf.String())
	}
}

func TestNotify_MeetsThreshold_SendsNotification(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf, notify.Options{MinDriftCount: 2})

	sent, err := n.Notify(makeResults("svc-a", "svc-b", "svc-c"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sent {
		t.Error("expected notification when drift meets threshold")
	}
	if !strings.Contains(buf.String(), "3 service(s)") {
		t.Errorf("expected count in output, got %q", buf.String())
	}
}

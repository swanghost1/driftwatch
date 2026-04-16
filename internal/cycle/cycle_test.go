package cycle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func drifted(service, field string) drift.Result {
	return drift.Result{Service: service, Field: field, Drifted: true}
}

func clean(service, field string) drift.Result {
	return drift.Result{Service: service, Field: field, Drifted: false}
}

func TestDetect_NoHistory_ReturnsNil(t *testing.T) {
	result := Detect(nil, DefaultOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestDetect_FewerRunsThanMinRuns_ReturnsNil(t *testing.T) {
	history := [][]drift.Result{
		{drifted("svc-a", "image")},
		{drifted("svc-a", "image")},
	}
	result := Detect(history, Options{MinRuns: 3})
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestDetect_PersistentDrift_Detected(t *testing.T) {
	history := [][]drift.Result{
		{drifted("svc-a", "image")},
		{drifted("svc-a", "image")},
		{drifted("svc-a", "image")},
	}
	results := Detect(history, Options{MinRuns: 3})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Service != "svc-a" || results[0].Field != "image" {
		t.Errorf("unexpected result: %+v", results[0])
	}
	if results[0].Consecutive != 3 {
		t.Errorf("expected 3 consecutive, got %d", results[0].Consecutive)
	}
}

func TestDetect_InterruptedDrift_NotCycle(t *testing.T) {
	history := [][]drift.Result{
		{drifted("svc-a", "image")},
		{clean("svc-a", "image")},
		{drifted("svc-a", "image")},
	}
	results := Detect(history, Options{MinRuns: 3})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestDetect_MixedServices_OnlyPersistentReturned(t *testing.T) {
	history := [][]drift.Result{
		{drifted("svc-a", "image"), drifted("svc-b", "replicas")},
		{drifted("svc-a", "image"), clean("svc-b", "replicas")},
		{drifted("svc-a", "image"), drifted("svc-b", "replicas")},
	}
	results := Detect(history, Options{MinRuns: 3})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Service != "svc-a" {
		t.Errorf("expected svc-a, got %s", results[0].Service)
	}
}

func TestWrite_NoResults_ShowsMessage(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, nil)
	if !strings.Contains(buf.String(), "no persistent") {
		t.Errorf("expected no-cycles message, got: %s", buf.String())
	}
}

func TestWrite_WithResults_ContainsServiceAndField(t *testing.T) {
	var buf bytes.Buffer
	Write(&buf, []Result{{Service: "svc-a", Field: "image", Consecutive: 5}})
	out := buf.String()
	if !strings.Contains(out, "svc-a") {
		t.Errorf("expected svc-a in output")
	}
	if !strings.Contains(out, "image") {
		t.Errorf("expected image in output")
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected count 5 in output")
	}
}

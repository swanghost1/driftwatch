package clip_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/clip"
	"github.com/example/driftwatch/internal/drift"
)

func makeResults(pairs ...string) []drift.Result {
	// pairs: service, field alternating
	var out []drift.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, drift.Result{
			Service: pairs[i],
			Field:   pairs[i+1],
			Drifted: true,
		})
	}
	return out
}

func TestApply_ZeroMax_ReturnsAll(t *testing.T) {
	results := makeResults("svc-a", "image", "svc-b", "image", "svc-c", "replicas")
	out := clip.Apply(results, clip.Options{MaxPerField: 0})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestApply_CapPerField_LimitsCorrectly(t *testing.T) {
	results := makeResults(
		"svc-a", "image",
		"svc-b", "image",
		"svc-c", "image",
		"svc-d", "replicas",
	)
	out := clip.Apply(results, clip.Options{MaxPerField: 2})
	imageCnt := 0
	for _, r := range out {
		if r.Field == "image" {
			imageCnt++
		}
	}
	if imageCnt != 2 {
		t.Fatalf("expected 2 image results, got %d", imageCnt)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 total results, got %d", len(out))
	}
}

func TestApply_PreservesOrder(t *testing.T) {
	results := makeResults("alpha", "image", "beta", "image", "gamma", "image")
	out := clip.Apply(results, clip.Options{MaxPerField: 2})
	if out[0].Service != "alpha" || out[1].Service != "beta" {
		t.Fatalf("order not preserved: %v", out)
	}
}

func TestApply_MultipleFields_IndependentCaps(t *testing.T) {
	results := makeResults(
		"a", "image", "b", "image", "c", "image",
		"d", "replicas", "e", "replicas", "f", "replicas",
	)
	out := clip.Apply(results, clip.Options{MaxPerField: 2})
	if len(out) != 4 {
		t.Fatalf("expected 4, got %d", len(out))
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := clip.Apply(nil, clip.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	results := makeResults("svc", "image")
	var buf bytes.Buffer
	clip.Write(&buf, results)
	if !strings.Contains(buf.String(), "FIELD") {
		t.Fatal("expected FIELD header")
	}
}

func TestWrite_ShowsFieldCounts(t *testing.T) {
	results := makeResults("a", "image", "b", "image", "c", "replicas")
	var buf bytes.Buffer
	clip.Write(&buf, results)
	if !strings.Contains(buf.String(), "image") {
		t.Fatal("expected image in output")
	}
}

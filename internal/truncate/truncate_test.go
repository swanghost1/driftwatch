package truncate_test

import (
	"bytes"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/truncate"
)

func makeResults(drifted, clean int) []drift.Result {
	var out []drift.Result
	for i := 0; i < drifted; i++ {
		out = append(out, drift.Result{Service: fmt.Sprintf("svc-drift-%d", i), HasDrift: true})
	}
	for i := 0; i < clean; i++ {
		out = append(out, drift.Result{Service: fmt.Sprintf("svc-ok-%d", i), HasDrift: false})
	}
	return out
}

func TestApply_NoMaxResults_ReturnsAll(t *testing.T) {
	results := makeResults(3, 3)
	got, dropped := truncate.Apply(results, truncate.Options{MaxResults: 0})
	if len(got) != 6 || dropped != 0 {
		t.Fatalf("expected 6/0, got %d/%d", len(got), dropped)
	}
}

func TestApply_BelowMax_ReturnsAll(t *testing.T) {
	results := makeResults(2, 2)
	got, dropped := truncate.Apply(results, truncate.Options{MaxResults: 10})
	if len(got) != 4 || dropped != 0 {
		t.Fatalf("expected 4/0, got %d/%d", len(got), dropped)
	}
}

func TestApply_TruncatesCorrectly(t *testing.T) {
	results := makeResults(0, 8)
	got, dropped := truncate.Apply(results, truncate.Options{MaxResults: 5})
	if len(got) != 5 {
		t.Fatalf("expected 5 results, got %d", len(got))
	}
	if dropped != 3 {
		t.Fatalf("expected 3 dropped, got %d", dropped)
	}
}

func TestApply_PrioritiseDrifted_DriftedKeptFirst(t *testing.T) {
	// 2 drifted + 4 clean; keep 3 → should retain both drifted entries
	results := makeResults(2, 4)
	got, dropped := truncate.Apply(results, truncate.Options{MaxResults: 3, PrioritiseDrifted: true})
	if len(got) != 3 || dropped != 3 {
		t.Fatalf("expected 3/3, got %d/%d", len(got), dropped)
	}
	driftCount := 0
	for _, r := range got {
		if r.HasDrift {
			driftCount++
		}
	}
	if driftCount != 2 {
		t.Fatalf("expected both drifted results retained, got %d", driftCount)
	}
}

func TestApply_DoesNotModifyOriginal(t *testing.T) {
	results := makeResults(1, 5)
	origFirst := results[0].Service
	truncate.Apply(results, truncate.Options{MaxResults: 2, PrioritiseDrifted: true})
	if results[0].Service != origFirst {
		t.Fatal("Apply modified the original slice")
	}
}

func TestWrite_NoneDropped(t *testing.T) {
	var buf bytes.Buffer
	truncate.Write(&buf, 5, 0)
	if !bytes.Contains(buf.Bytes(), []byte("retained")) {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestWrite_SomeDropped(t *testing.T) {
	var buf bytes.Buffer
	truncate.Write(&buf, 5, 3)
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("dropped")) {
		t.Fatalf("expected 'dropped' in output: %s", out)
	}
}

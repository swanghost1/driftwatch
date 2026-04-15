package ceiling_test

import (
	"bytes"
	"testing"

	"github.com/example/driftwatch/internal/ceiling"
	"github.com/example/driftwatch/internal/drift"
)

func makeResults(driftedCount, cleanCount int) []drift.Result {
	var out []drift.Result
	for i := 0; i < driftedCount; i++ {
		out = append(out, drift.Result{Service: fmt.Sprintf("svc-drifted-%d", i), Drifted: true})
	}
	for i := 0; i < cleanCount; i++ {
		out = append(out, drift.Result{Service: fmt.Sprintf("svc-clean-%d", i), Drifted: false})
	}
	return out
}

func TestApply_ZeroMax_ReturnsAll(t *testing.T) {
	results := makeResults(3, 3)
	opts := ceiling.DefaultOptions()
	got := ceiling.Apply(results, opts)
	if len(got) != 6 {
		t.Fatalf("expected 6 results, got %d", len(got))
	}
}

func TestApply_BelowMax_ReturnsAll(t *testing.T) {
	results := makeResults(2, 2)
	opts := ceiling.Options{MaxResults: 10, DriftedFirst: true}
	got := ceiling.Apply(results, opts)
	if len(got) != 4 {
		t.Fatalf("expected 4 results, got %d", len(got))
	}
}

func TestApply_TruncatesToMax(t *testing.T) {
	results := makeResults(4, 4)
	opts := ceiling.Options{MaxResults: 5, DriftedFirst: true}
	got := ceiling.Apply(results, opts)
	if len(got) != 5 {
		t.Fatalf("expected 5 results, got %d", len(got))
	}
}

func TestApply_DriftedFirst_DriftedResultsRetained(t *testing.T) {
	results := makeResults(3, 5)
	opts := ceiling.Options{MaxResults: 3, DriftedFirst: true}
	got := ceiling.Apply(results, opts)
	for _, r := range got {
		if !r.Drifted {
			t.Errorf("expected only drifted results when ceiling==driftedCount, got clean: %s", r.Service)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults(2, 2)
	orig := make([]drift.Result, len(results))
	copy(orig, results)
	opts := ceiling.Options{MaxResults: 2, DriftedFirst: true}
	ceiling.Apply(results, opts)
	for i, r := range results {
		if r.Service != orig[i].Service {
			t.Errorf("original slice mutated at index %d", i)
		}
	}
}

func TestWrite_Disabled(t *testing.T) {
	var buf bytes.Buffer
	ceiling.Write(&buf, 10, 10, ceiling.DefaultOptions())
	if !bytes.Contains(buf.Bytes(), []byte("disabled")) {
		t.Errorf("expected 'disabled' in output, got: %s", buf.String())
	}
}

func TestWrite_ShowsTrimmedCount(t *testing.T) {
	var buf bytes.Buffer
	opts := ceiling.Options{MaxResults: 5}
	ceiling.Write(&buf, 10, 5, opts)
	if !bytes.Contains(buf.Bytes(), []byte("trimmed=5")) {
		t.Errorf("expected trimmed=5 in output, got: %s", buf.String())
	}
}

package batch_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/batch"
	"github.com/driftwatch/internal/drift"
)

func makeResults(n int, driftedIndices ...int) []drift.Result {
	set := make(map[int]bool, len(driftedIndices))
	for _, i := range driftedIndices {
		set[i] = true
	}
	out := make([]drift.Result, n)
	for i := range out {
		out[i] = drift.Result{
			Service: fmt.Sprintf("svc-%d", i),
			Drifted: set[i],
		}
	}
	return out
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	batches := batch.Apply(nil, batch.DefaultOptions())
	if batches != nil {
		t.Fatalf("expected nil, got %v", batches)
	}
}

func TestApply_ZeroSize_ReturnsSingleBatch(t *testing.T) {
	results := makeResults(10)
	batches := batch.Apply(results, batch.Options{Size: 0})
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	if len(batches[0]) != 10 {
		t.Fatalf("expected 10 results in batch, got %d", len(batches[0]))
	}
}

func TestApply_ExactMultiple_CorrectBatchCount(t *testing.T) {
	results := makeResults(9)
	batches := batch.Apply(results, batch.Options{Size: 3})
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}
	for i, b := range batches {
		if len(b) != 3 {
			t.Errorf("batch %d: expected 3 results, got %d", i, len(b))
		}
	}
}

func TestApply_Remainder_LastBatchSmaller(t *testing.T) {
	results := makeResults(7)
	batches := batch.Apply(results, batch.Options{Size: 3})
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}
	if len(batches[2]) != 1 {
		t.Errorf("expected last batch to have 1 result, got %d", len(batches[2]))
	}
}

func TestApply_TotalResultsPreserved(t *testing.T) {
	results := makeResults(15)
	batches := batch.Apply(results, batch.Options{Size: 4})
	total := 0
	for _, b := range batches {
		total += len(b)
	}
	if total != 15 {
		t.Errorf("expected 15 total results across batches, got %d", total)
	}
}

func TestWrite_ContainsBatchCount(t *testing.T) {
	results := makeResults(6, 0, 2)
	batches := batch.Apply(results, batch.Options{Size: 3})
	var buf bytes.Buffer
	batch.Write(&buf, batches)
	out := buf.String()
	if !strings.Contains(out, "batches: 2") {
		t.Errorf("expected batch count header, got:\n%s", out)
	}
	if !strings.Contains(out, "drifted=1") {
		t.Errorf("expected drifted count in output, got:\n%s", out)
	}
}

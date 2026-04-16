package partition_test

import (
	"bytes"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/partition"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true},
		{Service: "api", Field: "replicas", Drifted: false},
		{Service: "worker", Field: "image", Drifted: true},
		{Service: "cache", Field: "image", Drifted: false},
	}
}

func TestApply_ByService_KeyCount(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "service"})
	if len(parts) != 3 {
		t.Fatalf("expected 3 partitions, got %d", len(parts))
	}
}

func TestApply_ByService_SortedAlphabetically(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "service"})
	keys := []string{parts[0].Key, parts[1].Key, parts[2].Key}
	expected := []string{"api", "cache", "worker"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q got %q", i, expected[i], k)
		}
	}
}

func TestApply_ByField_KeyCount(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "field"})
	if len(parts) != 2 {
		t.Fatalf("expected 2 partitions, got %d", len(parts))
	}
}

func TestApply_ByStatus_TwoBuckets(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "status"})
	if len(parts) != 2 {
		t.Fatalf("expected 2 partitions (ok/drifted), got %d", len(parts))
	}
}

func TestApply_ByStatus_DriftedCount(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "status"})
	for _, p := range parts {
		if p.Key == "drifted" && len(p.Results) != 2 {
			t.Errorf("expected 2 drifted results, got %d", len(p.Results))
		}
	}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	parts := partition.Apply(nil, partition.DefaultOptions())
	if parts != nil {
		t.Errorf("expected nil, got %v", parts)
	}
}

func TestApply_UnknownBy_FallsBackToService(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "unknown"})
	if len(parts) != 3 {
		t.Fatalf("expected 3 service partitions, got %d", len(parts))
	}
}

func TestWrite_ContainsKey(t *testing.T) {
	parts := partition.Apply(makeResults(), partition.Options{By: "service"})
	var buf bytes.Buffer
	partition.Write(&buf, parts)
	out := buf.String()
	for _, key := range []string{"api", "cache", "worker"} {
		if !bytes.Contains([]byte(out), []byte(key)) {
			t.Errorf("expected output to contain %q", key)
		}
	}
}

func TestWrite_EmptyPartitions_ShowsMessage(t *testing.T) {
	var buf bytes.Buffer
	partition.Write(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("no partitions")) {
		t.Error("expected 'no partitions' message")
	}
}

package bucketize_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/driftwatch/internal/bucketize"
)

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	buckets := []bucketize.Bucket{
		{Label: "2024-06-10", Start: time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), Total: 5, Drifted: 3},
	}
	if err := bucketize.WriteJSON(&buf, buckets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 element, got %d", len(out))
	}
}

func TestWriteJSON_ContainsDriftRate(t *testing.T) {
	var buf bytes.Buffer
	buckets := []bucketize.Bucket{
		{Label: "2024-06-10", Start: time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), Total: 4, Drifted: 2},
	}
	_ = bucketize.WriteJSON(&buf, buckets)
	var out []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)
	rate, ok := out[0]["drift_rate"].(float64)
	if !ok {
		t.Fatal("expected drift_rate field")
	}
	if rate != 0.5 {
		t.Errorf("expected drift_rate 0.5, got %f", rate)
	}
}

func TestWriteJSON_EmptyBuckets_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := bucketize.WriteJSON(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty array")
	}
}

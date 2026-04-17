package clusterby_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/clusterby"
	"github.com/example/driftwatch/internal/drift"
)

func makeResult(service, field string, drifted bool) drift.Result {
	return drift.Result{Service: service, Field: field, Drifted: drifted}
}

func TestApply_ByField_KeyCount(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-b", "image", false),
		makeResult("svc-c", "replicas", true),
	}
	clusters := clusterby.Apply(results, clusterby.KeyField)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}

func TestApply_ByStatus_TwoBuckets(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-b", "image", false),
		makeResult("svc-c", "replicas", true),
	}
	clusters := clusterby.Apply(results, clusterby.KeyStatus)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
	names := map[string]bool{}
	for _, c := range clusters {
		names[c.Name] = true
	}
	if !names["drifted"] || !names["ok"] {
		t.Errorf("expected drifted and ok clusters, got %v", names)
	}
}

func TestApply_ByPrefix_GroupsByFirstSegment(t *testing.T) {
	results := []drift.Result{
		makeResult("api-gateway", "image", true),
		makeResult("api-auth", "replicas", false),
		makeResult("worker-email", "image", true),
	}
	clusters := clusterby.Apply(results, clusterby.KeyPrefix)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	result := clusterby.Apply(nil, clusterby.KeyField)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestApply_SortedAlphabetically(t *testing.T) {
	results := []drift.Result{
		makeResult("svc", "replicas", false),
		makeResult("svc", "image", true),
		makeResult("svc", "env", true),
	}
	clusters := clusterby.Apply(results, clusterby.KeyField)
	for i := 1; i < len(clusters); i++ {
		if clusters[i].Name < clusters[i-1].Name {
			t.Errorf("clusters not sorted: %s before %s", clusters[i-1].Name, clusters[i].Name)
		}
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", true)}
	clusters := clusterby.Apply(results, clusterby.KeyField)
	var buf bytes.Buffer
	clusterby.Write(&buf, clusters)
	if !strings.Contains(buf.String(), "CLUSTER") {
		t.Errorf("expected CLUSTER header in output")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-b", "replicas", false),
	}
	clusters := clusterby.Apply(results, clusterby.KeyField)
	var buf bytes.Buffer
	if err := clusterby.WriteJSON(&buf, clusters); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out []clusterby.Cluster
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Errorf("invalid JSON: %v", err)
	}
}

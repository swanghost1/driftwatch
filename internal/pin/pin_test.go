package pin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true, Declared: "nginx:1.25", Live: "nginx:1.24"},
		{Service: "api", Field: "replicas", Drifted: true, Declared: "3", Live: "2"},
		{Service: "worker", Field: "image", Drifted: false, Declared: "alpine:3", Live: "alpine:3"},
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	s := &Store{}
	out := Apply(makeResults(), s)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
	if !out[0].Drifted {
		t.Error("expected first result to remain drifted")
	}
}

func TestApply_MatchingRule_PinsResult(t *testing.T) {
	s := &Store{Rules: []Rule{{Service: "api", Field: "image", Reason: "intentional"}}}
	out := Apply(makeResults(), s)
	if out[0].Drifted {
		t.Error("expected pinned result to be marked not drifted")
	}
	if out[1].Drifted == false {
		t.Error("expected non-pinned drift to remain drifted")
	}
}

func TestApply_ExpiredRule_DoesNotPin(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	s := &Store{Rules: []Rule{{Service: "api", Field: "image", ExpiresAt: &expired}}}
	out := Apply(makeResults(), s)
	if !out[0].Drifted {
		t.Error("expired pin should not suppress drift")
	}
}

func TestApply_WildcardService_PinsAll(t *testing.T) {
	s := &Store{Rules: []Rule{{Service: "*", Field: "image"}}}
	out := Apply(makeResults(), s)
	// api/image and worker/image — worker is not drifted so stays unchanged
	if out[0].Drifted {
		t.Error("expected api/image to be pinned")
	}
}

func TestApply_EmptyField_PinsAllFieldsForService(t *testing.T) {
	s := &Store{Rules: []Rule{{Service: "api", Field: ""}}}
	out := Apply(makeResults(), s)
	if out[0].Drifted || out[1].Drifted {
		t.Error("expected all api drifts to be pinned")
	}
}

func TestLoadStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	s, err := LoadStore(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 0 {
		t.Errorf("expected empty store, got %d rules", len(s.Rules))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	s := &Store{path: path, Rules: []Rule{{Service: "api", Field: "image", Reason: "ok"}}}
	if err := SaveStore(s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadStore(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].Service != "api" {
		t.Errorf("unexpected rules: %+v", loaded.Rules)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "dir", "pins.json")
	s := &Store{path: path, Rules: []Rule{{Service: "svc", Field: "replicas"}}}
	if err := SaveStore(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	var loaded Store
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
}

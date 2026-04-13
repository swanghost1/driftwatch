package mute_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/mute"
)

var (
	now    = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	future = now.Add(24 * time.Hour)
	past   = now.Add(-24 * time.Hour)
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Drifted: true},
		{Service: "api", Field: "replicas", Drifted: true},
		{Service: "worker", Field: "image", Drifted: true},
	}
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	results := makeResults()
	out := mute.Apply(results, &mute.Store{}, now)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_MatchingRule_MutesResult(t *testing.T) {
	s := &mute.Store{
		Rules: []mute.Rule{
			{Service: "api", Field: "image", ExpiresAt: future},
		},
	}
	out := mute.Apply(makeResults(), s, now)
	for _, r := range out {
		if r.Service == "api" && r.Field == "image" {
			t.Fatal("muted result should have been removed")
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_ExpiredRule_DoesNotMute(t *testing.T) {
	s := &mute.Store{
		Rules: []mute.Rule{
			{Service: "api", Field: "image", ExpiresAt: past},
		},
	}
	out := mute.Apply(makeResults(), s, now)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestApply_EmptyField_MutesAllFieldsForService(t *testing.T) {
	s := &mute.Store{
		Rules: []mute.Rule{
			{Service: "api", Field: "", ExpiresAt: future},
		},
	}
	out := mute.Apply(makeResults(), s, now)
	for _, r := range out {
		if r.Service == "api" {
			t.Fatalf("all api results should be muted, found field=%s", r.Field)
		}
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	s := &mute.Store{
		Rules: []mute.Rule{
			{Service: "API", Field: "IMAGE", ExpiresAt: future},
		},
	}
	out := mute.Apply(makeResults(), s, now)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestLoadStore_FileNotFound_ReturnsEmpty(t *testing.T) {
	s, err := mute.LoadStore(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mutes", "store.json")
	s := &mute.Store{
		Rules: []mute.Rule{
			{Service: "api", Field: "image", ExpiresAt: future, Reason: "planned deploy"},
		},
	}
	if err := mute.SaveStore(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := mute.LoadStore(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(loaded.Rules))
	}
	if loaded.Rules[0].Reason != "planned deploy" {
		t.Errorf("reason mismatch: %q", loaded.Rules[0].Reason)
	}
}

func TestSave_CreatesIntermediateDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c", "store.json")
	if err := mute.SaveStore(path, &mute.Store{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
	var s mute.Store
	data, _ := os.ReadFile(path)
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
}

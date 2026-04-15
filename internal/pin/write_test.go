package pin

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func activeStore() *Store {
	future := time.Now().Add(24 * time.Hour)
	return &Store{
		Rules: []Rule{
			{Service: "api", Field: "image", Reason: "planned migration", PinnedAt: time.Now(), ExpiresAt: &future},
			{Service: "worker", Field: "", Reason: "maintenance"},
		},
	}
}

func TestWriteText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteText(&buf, activeStore()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Error("expected header row in output")
	}
}

func TestWriteText_ShowsAllFieldFallback(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, activeStore())
	if !strings.Contains(buf.String(), "(all)") {
		t.Error("expected (all) for empty field pin")
	}
}

func TestWriteText_NoActiveRules_ShowsMessage(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	s := &Store{Rules: []Rule{{Service: "api", Field: "image", ExpiresAt: &expired}}}
	var buf bytes.Buffer
	WriteText(&buf, s)
	if !strings.Contains(buf.String(), "No active") {
		t.Error("expected no-active message")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, activeStore()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []Rule
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 rules, got %d", len(out))
	}
}

func TestWriteJSON_ExcludesExpired(t *testing.T) {
	expired := time.Now().Add(-time.Minute)
	s := &Store{Rules: []Rule{
		{Service: "api", Field: "image"},
		{Service: "old", Field: "replicas", ExpiresAt: &expired},
	}}
	var buf bytes.Buffer
	WriteJSON(&buf, s)
	var out []Rule
	json.Unmarshal(buf.Bytes(), &out)
	if len(out) != 1 {
		t.Errorf("expected 1 active rule, got %d", len(out))
	}
}

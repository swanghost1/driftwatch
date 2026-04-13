package retention

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWrite_NoPruned_ContainsPolicyInfo(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		Pruned: nil,
		Policy: DefaultPolicy(),
		RanAt:  time.Now(),
	}
	if err := Write(&buf, r); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "max-entries=100") {
		t.Errorf("expected max-entries in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Deleted:") {
		t.Errorf("expected Deleted: in output, got:\n%s", out)
	}
}

func TestWrite_WithPruned_ListsFiles(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		Pruned: []string{"/tmp/old_entry.json", "/tmp/older_entry.json"},
		Policy: Policy{MaxAge: 24 * 3600 * 1e9, MaxEntries: 50},
		RanAt:  time.Now(),
	}
	if err := Write(&buf, r); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "old_entry.json") {
		t.Errorf("expected pruned filename in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Removed files:") {
		t.Errorf("expected Removed files: section, got:\n%s", out)
	}
}

func TestWrite_ZeroPruned_ShowsZero(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		Pruned: []string{},
		Policy: DefaultPolicy(),
		RanAt:  time.Now(),
	}
	if err := Write(&buf, r); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 file(s)") {
		t.Errorf("expected '0 file(s)' in output, got:\n%s", out)
	}
}

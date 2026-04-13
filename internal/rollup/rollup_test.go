package rollup_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/rollup"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "frontend", Drifted: true},
		{Service: "backend", Drifted: false},
		{Service: "worker", Drifted: true},
		{Service: "cache", Drifted: false},
		{Service: "gateway", Drifted: true},
	}
}

func prefixKey(service string) string {
	switch service {
	case "frontend", "backend":
		return "web"
	case "worker", "cache":
		return "infra"
	default:
		return ""
	}
}

func TestByGroup_GroupCount(t *testing.T) {
	groups := rollup.ByGroup(makeResults(), prefixKey)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestByGroup_Totals(t *testing.T) {
	groups := rollup.ByGroup(makeResults(), prefixKey)
	byName := map[string]rollup.Group{}
	for _, g := range groups {
		byName[g.Name] = g
	}

	if byName["web"].Total != 2 {
		t.Errorf("web total: want 2, got %d", byName["web"].Total)
	}
	if byName["infra"].Drifted != 1 {
		t.Errorf("infra drifted: want 1, got %d", byName["infra"].Drifted)
	}
	if byName["(untagged)"].Total != 1 {
		t.Errorf("untagged total: want 1, got %d", byName["(untagged)"].Total)
	}
}

func TestByGroup_DriftRate(t *testing.T) {
	groups := rollup.ByGroup(makeResults(), prefixKey)
	byName := map[string]rollup.Group{}
	for _, g := range groups {
		byName[g.Name] = g
	}

	if byName["web"].DriftRate != 50.0 {
		t.Errorf("web drift rate: want 50.0, got %.1f", byName["web"].DriftRate)
	}
}

func TestByGroup_EmptyResults(t *testing.T) {
	groups := rollup.ByGroup(nil, prefixKey)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestWrite_ContainsHeaders(t *testing.T) {
	groups := rollup.ByGroup(makeResults(), prefixKey)
	var buf bytes.Buffer
	if err := rollup.Write(&buf, groups); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"GROUP", "TOTAL", "DRIFTED", "CLEAN", "DRIFT%"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
}

func TestWrite_ContainsGroupNames(t *testing.T) {
	groups := rollup.ByGroup(makeResults(), prefixKey)
	var buf bytes.Buffer
	_ = rollup.Write(&buf, groups)
	out := buf.String()
	for _, name := range []string{"web", "infra", "(untagged)"} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing group %q", name)
		}
	}
}

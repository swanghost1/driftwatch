package rollup_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/rollup"
)

// TestRoundTrip_AllDrifted verifies that a group containing only drifted
// services reports a 100 % drift rate and zero clean services.
func TestRoundTrip_AllDrifted(t *testing.T) {
	results := []drift.Result{
		{Service: "svc-a", Drifted: true},
		{Service: "svc-b", Drifted: true},
	}
	groups := rollup.ByGroup(results, func(_ string) string { return "prod" })
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	g := groups[0]
	if g.Clean != 0 {
		t.Errorf("clean: want 0, got %d", g.Clean)
	}
	if g.DriftRate != 100.0 {
		t.Errorf("drift rate: want 100.0, got %.1f", g.DriftRate)
	}
}

// TestRoundTrip_NoDrift verifies that a fully-clean group reports 0 % drift.
func TestRoundTrip_NoDrift(t *testing.T) {
	results := []drift.Result{
		{Service: "svc-a", Drifted: false},
		{Service: "svc-b", Drifted: false},
	}
	groups := rollup.ByGroup(results, func(_ string) string { return "staging" })
	if groups[0].DriftRate != 0.0 {
		t.Errorf("drift rate: want 0.0, got %.1f", groups[0].DriftRate)
	}
}

// TestWrite_SortedOutput verifies groups are emitted in alphabetical order.
func TestWrite_SortedOutput(t *testing.T) {
	results := []drift.Result{
		{Service: "z-svc", Drifted: false},
		{Service: "a-svc", Drifted: true},
		{Service: "m-svc", Drifted: false},
	}
	groups := rollup.ByGroup(results, func(s string) string { return s })
	var buf bytes.Buffer
	_ = rollup.Write(&buf, groups)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// lines[0] = header, lines[1] = separator, lines[2..] = data
	if len(lines) < 5 {
		t.Fatalf("expected at least 5 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[2]), "a-svc") {
		t.Errorf("first data row should be a-svc, got: %s", lines[2])
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[4]), "z-svc") {
		t.Errorf("last data row should be z-svc, got: %s", lines[4])
	}
}

package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

// TestRoundTrip_WarnThenError verifies that a mixed result set triggers the
// correct rules when two severity levels are configured.
func TestRoundTrip_WarnThenError(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Drifted: true},
		{Service: "worker", Drifted: true},
		{Service: "cache", Drifted: false},
	}

	rules := []alert.Rule{
		{MinDrifted: 1, Level: alert.LevelWarn, Label: "any-drift"},
		{MinDrifted: 3, Level: alert.LevelError, Label: "all-drift"},
	}

	alerts := alert.Evaluate(results, rules)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert (warn only), got %d", len(alerts))
	}
	if alerts[0].Rule.Level != alert.LevelWarn {
		t.Errorf("expected warn level, got %s", alerts[0].Rule.Level)
	}

	var buf bytes.Buffer
	alert.Write(&buf, alerts)
	out := buf.String()

	if !strings.Contains(out, "api") || !strings.Contains(out, "worker") {
		t.Errorf("expected drifted service names in output: %s", out)
	}
	if strings.Contains(out, "cache") {
		t.Errorf("clean service should not appear in alert output: %s", out)
	}
}

// TestRoundTrip_AllDrifted_BothRulesTriggered ensures all matching rules fire.
func TestRoundTrip_AllDrifted_BothRulesTriggered(t *testing.T) {
	results := []drift.Result{
		{Service: "api", Drifted: true},
		{Service: "worker", Drifted: true},
		{Service: "cache", Drifted: true},
	}

	rules := []alert.Rule{
		{MinDrifted: 1, Level: alert.LevelWarn, Label: "any-drift"},
		{MinDrifted: 3, Level: alert.LevelError, Label: "all-drift"},
	}

	alerts := alert.Evaluate(results, rules)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}

	levels := map[alert.Level]bool{}
	for _, a := range alerts {
		levels[a.Rule.Level] = true
	}
	if !levels[alert.LevelWarn] || !levels[alert.LevelError] {
		t.Errorf("expected both warn and error levels, got: %v", levels)
	}
}

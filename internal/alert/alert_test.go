package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func makeResults(services []string, drifted []bool) []drift.Result {
	var out []drift.Result
	for i, name := range services {
		out = append(out, drift.Result{Service: name, Drifted: drifted[i]})
	}
	return out
}

func TestEvaluate_NoDrift_NoAlerts(t *testing.T) {
	results := makeResults([]string{"svc-a", "svc-b"}, []bool{false, false})
	rules := []alert.Rule{{MinDrifted: 1, Level: alert.LevelWarn, Label: "any-drift"}}

	alerts := alert.Evaluate(results, rules)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}

func TestEvaluate_BelowThreshold_NoAlert(t *testing.T) {
	results := makeResults([]string{"svc-a", "svc-b"}, []bool{true, false})
	rules := []alert.Rule{{MinDrifted: 2, Level: alert.LevelError, Label: "many-drift"}}

	alerts := alert.Evaluate(results, rules)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}

func TestEvaluate_MeetsThreshold_RaisesAlert(t *testing.T) {
	results := makeResults([]string{"svc-a", "svc-b"}, []bool{true, true})
	rules := []alert.Rule{{MinDrifted: 2, Level: alert.LevelError, Label: "many-drift"}}

	alerts := alert.Evaluate(results, rules)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Rule.Level != alert.LevelError {
		t.Errorf("expected level error, got %s", alerts[0].Rule.Level)
	}
	if len(alerts[0].Drifted) != 2 {
		t.Errorf("expected 2 drifted services, got %d", len(alerts[0].Drifted))
	}
}

func TestEvaluate_MultipleRules_AllMatching(t *testing.T) {
	results := makeResults([]string{"svc-a", "svc-b", "svc-c"}, []bool{true, true, true})
	rules := []alert.Rule{
		{MinDrifted: 1, Level: alert.LevelWarn, Label: "warn-rule"},
		{MinDrifted: 3, Level: alert.LevelError, Label: "error-rule"},
	}

	alerts := alert.Evaluate(results, rules)
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestWrite_NoAlerts_PrintsNone(t *testing.T) {
	var buf bytes.Buffer
	alert.Write(&buf, nil)
	if !strings.Contains(buf.String(), "no alerts") {
		t.Errorf("expected 'no alerts' message, got: %s", buf.String())
	}
}

func TestWrite_WithAlerts_ContainsLevel(t *testing.T) {
	results := makeResults([]string{"svc-a"}, []bool{true})
	rules := []alert.Rule{{MinDrifted: 1, Level: alert.LevelWarn, Label: "test-rule"}}
	alerts := alert.Evaluate(results, rules)

	var buf bytes.Buffer
	alert.Write(&buf, alerts)
	out := buf.String()

	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "test-rule") {
		t.Errorf("expected rule label in output, got: %s", out)
	}
	if !strings.Contains(out, "svc-a") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}

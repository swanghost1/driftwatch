// Package alert provides threshold-based alerting for drift results,
// routing notifications to configured channels when drift counts exceed limits.
package alert

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Rule defines the conditions under which an alert is raised.
type Rule struct {
	// MinDrifted is the minimum number of drifted services to trigger the alert.
	MinDrifted int
	// Level is the severity assigned to the alert when triggered.
	Level Level
	// Label is an optional human-readable name for the rule.
	Label string
}

// Alert is raised when a Rule is matched.
type Alert struct {
	Rule    Rule
	Drifted []string
	Message string
}

// Evaluate checks results against a set of rules and returns any triggered alerts.
// Rules are evaluated in order; all matching rules produce an alert.
func Evaluate(results []drift.Result, rules []Rule) []Alert {
	var alerts []Alert

	drifted := driftedNames(results)
	count := len(drifted)

	for _, rule := range rules {
		if count >= rule.MinDrifted {
			msg := fmt.Sprintf("%s: %d service(s) drifted (%s)",
				rule.Level,
				count,
				strings.Join(drifted, ", "),
			)
			alerts = append(alerts, Alert{
				Rule:    rule,
				Drifted: drifted,
				Message: msg,
			})
		}
	}

	return alerts
}

// Write renders alerts as human-readable lines to w.
func Write(w io.Writer, alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "no alerts triggered")
		return
	}
	for _, a := range alerts {
		label := a.Rule.Label
		if label == "" {
			label = "unnamed"
		}
		fmt.Fprintf(w, "[%s] rule=%q %s\n", strings.ToUpper(string(a.Rule.Level)), label, a.Message)
	}
}

func driftedNames(results []drift.Result) []string {
	var names []string
	for _, r := range results {
		if r.Drifted {
			names = append(names, r.Service)
		}
	}
	return names
}

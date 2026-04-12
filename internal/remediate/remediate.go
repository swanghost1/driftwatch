// Package remediate suggests corrective actions for detected config drift.
package remediate

import (
	"fmt"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Suggestion holds a human-readable remediation hint for a single drift result.
type Suggestion struct {
	Service    string
	Field      string
	Expected   string
	Actual     string
	Action     string
}

// Suggest returns a slice of Suggestions for all drifted results.
func Suggest(results []drift.Result) []Suggestion {
	var suggestions []Suggestion
	for _, r := range results {
		if !r.Drifted {
			continue
		}
		for _, d := range r.Diffs {
			suggestions = append(suggestions, Suggestion{
				Service:  r.Service,
				Field:    d.Field,
				Expected: fmt.Sprintf("%v", d.Expected),
				Actual:   fmt.Sprintf("%v", d.Actual),
				Action:   buildAction(r.Service, d.Field, d.Expected),
			})
		}
	}
	return suggestions
}

// buildAction constructs a kubectl-style command hint for a given field drift.
func buildAction(service, field string, expected interface{}) string {
	switch strings.ToLower(field) {
	case "image":
		return fmt.Sprintf(
			"kubectl set image deployment/%s *=%v",
			service, expected,
		)
	case "replicas":
		return fmt.Sprintf(
			"kubectl scale deployment/%s --replicas=%v",
			service, expected,
		)
	default:
		if strings.HasPrefix(strings.ToLower(field), "env:") {
			envKey := strings.TrimPrefix(strings.ToLower(field), "env:")
			return fmt.Sprintf(
				"kubectl set env deployment/%s %s=%v",
				service, strings.ToUpper(envKey), expected,
			)
		}
		return fmt.Sprintf(
			"# manually reconcile field '%s' on deployment/%s to: %v",
			field, service, expected,
		)
	}
}

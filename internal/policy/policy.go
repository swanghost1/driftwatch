// Package policy evaluates drift results against user-defined rules
// and determines whether a run should be considered a failure.
package policy

import (
	"fmt"

	"github.com/example/driftwatch/internal/drift"
)

// Rule defines a single policy constraint.
type Rule struct {
	// MaxDriftedServices is the maximum number of drifted services allowed
	// before the policy is violated. Zero means any drift is a violation.
	MaxDriftedServices int `yaml:"max_drifted_services"`

	// BlockedFields lists field names whose drift always causes a violation,
	// regardless of MaxDriftedServices.
	BlockedFields []string `yaml:"blocked_fields"`
}

// Violation describes a single policy breach.
type Violation struct {
	Service string
	Field   string
	Reason  string
}

// Evaluate checks results against the given Rule and returns any violations.
// A non-nil, non-empty slice means the policy was breached.
func Evaluate(results []drift.Result, rule Rule) []Violation {
	var violations []Violation

	blockedSet := make(map[string]bool, len(rule.BlockedFields))
	for _, f := range rule.BlockedFields {
		blockedSet[f] = true
	}

	driftedCount := 0
	for _, r := range results {
		if !r.Drifted {
			continue
		}
		driftedCount++

		for _, d := range r.Diffs {
			if blockedSet[d.Field] {
				violations = append(violations, Violation{
					Service: r.Service,
					Field:   d.Field,
					Reason:  fmt.Sprintf("field %q is blocked by policy", d.Field),
				})
			}
		}
	}

	if driftedCount > rule.MaxDriftedServices {
		violations = append(violations, Violation{
			Reason: fmt.Sprintf(
				"%d drifted service(s) exceed allowed maximum of %d",
				driftedCount, rule.MaxDriftedServices,
			),
		})
	}

	return violations
}

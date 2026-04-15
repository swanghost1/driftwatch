// Package clamp enforces minimum and maximum replica bounds on drift results,
// flagging any service whose live replica count falls outside the declared range.
package clamp

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Options controls the clamping behaviour.
type Options struct {
	// MinReplicas is the lowest acceptable replica count (0 = no lower bound).
	MinReplicas int
	// MaxReplicas is the highest acceptable replica count (0 = no upper bound).
	MaxReplicas int
}

// Violation describes a single out-of-bounds replica count.
type Violation struct {
	Service     string
	LiveValue   int
	Declared    int
	Min         int
	Max         int
	Description string
}

// Apply inspects each drift result and returns violations where the live
// replica count is outside [opts.MinReplicas, opts.MaxReplicas].
// Results that are not replica-related are skipped.
func Apply(results []drift.Result, opts Options) []Violation {
	var violations []Violation
	for _, r := range results {
		if r.Field != "replicas" {
			continue
		}
		var live, declared int
		if _, err := fmt.Sscanf(r.LiveValue, "%d", &live); err != nil {
			continue
		}
		if _, err := fmt.Sscanf(r.WantValue, "%d", &declared); err != nil {
			continue
		}
		var desc string
		if opts.MinReplicas > 0 && live < opts.MinReplicas {
			desc = fmt.Sprintf("live replicas %d below minimum %d", live, opts.MinReplicas)
		} else if opts.MaxReplicas > 0 && live > opts.MaxReplicas {
			desc = fmt.Sprintf("live replicas %d above maximum %d", live, opts.MaxReplicas)
		}
		if desc != "" {
			violations = append(violations, Violation{
				Service:     r.Service,
				LiveValue:   live,
				Declared:    declared,
				Min:         opts.MinReplicas,
				Max:         opts.MaxReplicas,
				Description: desc,
			})
		}
	}
	return violations
}

// Write renders violations as a human-readable table to w.
func Write(w io.Writer, violations []Violation) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tLIVE\tDECLARED\tMIN\tMAX\tDESCRIPTION")
	for _, v := range violations {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%s\n",
			v.Service, v.LiveValue, v.Declared, v.Min, v.Max, v.Description)
	}
	tw.Flush()
}

// Package quota enforces per-service drift count limits, flagging services
// that exceed a configured maximum number of drifted fields in a single run.
package quota

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/example/driftwatch/internal/drift"
)

// Options configures quota enforcement.
type Options struct {
	// MaxDriftedFields is the maximum number of drifted fields allowed per
	// service. Zero means no limit is enforced.
	MaxDriftedFields int
}

// Violation describes a service that has exceeded its quota.
type Violation struct {
	Service string
	DriftedFields int
	Limit int
}

// Result holds the outcome of a quota evaluation.
type Result struct {
	Violations []Violation
}

// Exceeded reports whether any quota was exceeded.
func (r Result) Exceeded() bool { return len(r.Violations) > 0 }

// Evaluate checks each service in results against opts and returns any
// violations. Results are unchanged; this is a read-only inspection.
func Evaluate(results []drift.Result, opts Options) Result {
	if opts.MaxDriftedFields <= 0 {
		return Result{}
	}

	counts := make(map[string]int)
	for _, r := range results {
		if r.Drifted {
			counts[r.Service]++
		}
	}

	var violations []Violation
	for svc, n := range counts {
		if n > opts.MaxDriftedFields {
			violations = append(violations, Violation{
				Service:       svc,
				DriftedFields: n,
				Limit:         opts.MaxDriftedFields,
			})
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Service < violations[j].Service
	})

	return Result{Violations: violations}
}

// Write renders the quota result as a human-readable table to w.
func Write(w io.Writer, r Result) error {
	if !r.Exceeded() {
		_, err := fmt.Fprintln(w, "quota: no violations")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tDRIFTED FIELDS\tLIMIT")
	fmt.Fprintln(tw, "-------\t--------------\t-----")
	for _, v := range r.Violations {
		fmt.Fprintf(tw, "%s\t%d\t%d\n", v.Service, v.DriftedFields, v.Limit)
	}
	return tw.Flush()
}

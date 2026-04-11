// Package filter provides utilities for narrowing drift detection
// results based on service name patterns and drift severity.
package filter

import (
	"strings"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Options holds the filtering criteria applied to drift results.
type Options struct {
	// Services is an optional list of service name substrings to include.
	// An empty slice means "include all services".
	Services []string

	// OnlyDrifted, when true, drops results that have no drift items.
	OnlyDrifted bool
}

// Apply returns a filtered copy of results according to opts.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if !matchesService(r.ServiceName, opts.Services) {
			continue
		}
		if opts.OnlyDrifted && len(r.Drifts) == 0 {
			continue
		}
		out = append(out, r)
	}
	return out
}

// matchesService returns true when name matches any of the provided patterns
// (case-insensitive substring match). If patterns is empty, every name matches.
func matchesService(name string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	lower := strings.ToLower(name)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Package groupfilter provides filtering of drift results by group prefix,
// allowing operators to scope checks to a logical subset of services.
package groupfilter

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls which groups are included or excluded.
type Options struct {
	// IncludeGroups limits results to services whose group matches any entry.
	// Matching is case-insensitive prefix match on the service name.
	IncludeGroups []string

	// ExcludeGroups removes results whose group matches any entry.
	ExcludeGroups []string
}

// Apply returns a filtered copy of results according to opts.
// If IncludeGroups is non-empty, only matching results are kept.
// ExcludeGroups is applied after inclusion filtering.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if len(opts.IncludeGroups) > 0 && !matchesAny(r.Service, opts.IncludeGroups) {
			continue
		}
		if matchesAny(r.Service, opts.ExcludeGroups) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// Groups returns the distinct group prefixes present in results.
// A group is defined as the portion of the service name before the first '/'.
// Services without a '/' are placed in the "default" group.
func Groups(results []drift.Result) []string {
	seen := map[string]struct{}{}
	for _, r := range results {
		seen[groupOf(r.Service)] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for g := range seen {
		out = append(out, g)
	}
	return out
}

func groupOf(service string) string {
	if idx := strings.Index(service, "/"); idx >= 0 {
		return service[:idx]
	}
	return "default"
}

func matchesAny(service string, groups []string) bool {
	g := strings.ToLower(groupOf(service))
	for _, candidate := range groups {
		if strings.ToLower(candidate) == g {
			return true
		}
	}
	return false
}

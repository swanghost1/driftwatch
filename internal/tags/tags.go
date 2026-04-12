// Package tags provides utilities for filtering and grouping drift results
// by user-defined service tags declared in the config.
package tags

import (
	"sort"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// GroupByTag partitions results into groups keyed by the given tag name.
// Services that do not carry the tag are placed under the "(untagged)" key.
func GroupByTag(results []drift.Result, tag string) map[string][]drift.Result {
	groups := make(map[string][]drift.Result)
	for _, r := range results {
		val, ok := r.Tags[tag]
		if !ok || strings.TrimSpace(val) == "" {
			val = "(untagged)"
		}
		groups[val] = append(groups[val], r)
	}
	return groups
}

// Keys returns the sorted group keys produced by GroupByTag.
func Keys(groups map[string][]drift.Result) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FilterByTag returns only the results whose tag value matches val (case-insensitive).
// If val is empty, all results are returned unchanged.
func FilterByTag(results []drift.Result, tag, val string) []drift.Result {
	if strings.TrimSpace(val) == "" {
		return results
	}
	val = strings.ToLower(val)
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if strings.ToLower(r.Tags[tag]) == val {
			out = append(out, r)
		}
	}
	return out
}

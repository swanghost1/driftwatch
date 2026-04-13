// Package dedupe provides utilities for removing duplicate drift results
// across multiple detection runs, keyed by service name and field.
package dedupe

import (
	"fmt"

	"github.com/example/driftwatch/internal/drift"
)

// Key uniquely identifies a drift result by service and field.
type Key struct {
	Service string
	Field   string
}

// keyFor builds a deduplication key from a DetectResult.
func keyFor(r drift.DetectResult) Key {
	return Key{Service: r.Service, Field: r.Field}
}

// String returns a human-readable representation of the key.
func (k Key) String() string {
	return fmt.Sprintf("%s/%s", k.Service, k.Field)
}

// Apply removes duplicate DetectResult entries from results, keeping the first
// occurrence of each (service, field) pair. Results with no drift are always
// retained but are also deduplicated by service name.
func Apply(results []drift.DetectResult) []drift.DetectResult {
	seen := make(map[Key]struct{}, len(results))
	out := make([]drift.DetectResult, 0, len(results))

	for _, r := range results {
		k := keyFor(r)
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, r)
	}

	return out
}

// Count returns the number of unique (service, field) pairs in results.
func Count(results []drift.DetectResult) int {
	seen := make(map[Key]struct{}, len(results))
	for _, r := range results {
		seen[keyFor(r)] = struct{}{}
	}
	return len(seen)
}

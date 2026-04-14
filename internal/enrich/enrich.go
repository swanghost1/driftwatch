// Package enrich attaches contextual metadata to drift results,
// such as the time the drift was first detected and a human-readable
// description of the drifted field.
package enrich

import (
	"fmt"
	"strings"
	"time"
)

// Result mirrors the drift result shape used across the project.
type Result struct {
	Service   string
	Field     string
	Expected  string
	Actual    string
	Drifted   bool
	DetectedAt time.Time
	Description string
}

// Options controls which enrichment steps are applied.
type Options struct {
	// DetectedAt is stamped onto every result that has a zero time.
	DetectedAt time.Time
	// Describe generates a human-readable description for drifted fields.
	// When nil the built-in describer is used.
	Describe func(field, expected, actual string) string
}

// Apply enriches a slice of results in-place and returns it.
func Apply(results []Result, opts Options) []Result {
	ts := opts.DetectedAt
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	describe := opts.Describe
	if describe == nil {
		describe = defaultDescription
	}
	for i := range results {
		if results[i].DetectedAt.IsZero() {
			results[i].DetectedAt = ts
		}
		if results[i].Drifted && results[i].Description == "" {
			results[i].Description = describe(
				results[i].Field,
				results[i].Expected,
				results[i].Actual,
			)
		}
	}
	return results
}

func defaultDescription(field, expected, actual string) string {
	return fmt.Sprintf("%s: expected %q but got %q",
		strings.ToLower(field), expected, actual)
}

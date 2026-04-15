// Package mask provides field-level value masking for drift results.
// It replaces live and declared values for nominated fields with a
// fixed placeholder so that sensitive data is never written to reports
// while still indicating that a drift exists.
package mask

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

const placeholder = "***"

// Options controls which fields are masked.
type Options struct {
	// Fields is the list of field names whose values should be masked.
	// Matching is case-insensitive.
	Fields []string
	// Disabled skips all masking when true.
	Disabled bool
}

// DefaultOptions returns an Options with a sensible set of masked fields.
func DefaultOptions() Options {
	return Options{
		Fields: []string{"password", "secret", "token", "api_key", "apikey"},
	}
}

// Apply returns a copy of results with sensitive field values replaced by
// the mask placeholder. Results for non-masked fields are returned unchanged.
func Apply(results []drift.Result, opts Options) []drift.Result {
	if opts.Disabled || len(opts.Fields) == 0 {
		return results
	}

	out := make([]drift.Result, len(results))
	for i, r := range results {
		if isMasked(r.Field, opts.Fields) {
			r.Live = placeholder
			r.Declared = placeholder
		}
		out[i] = r
	}
	return out
}

// isMasked reports whether field matches any entry in the mask list.
func isMasked(field string, fields []string) bool {
	lower := strings.ToLower(field)
	for _, f := range fields {
		if strings.ToLower(f) == lower {
			return true
		}
	}
	return false
}

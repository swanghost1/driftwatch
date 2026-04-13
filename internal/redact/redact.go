// Package redact provides utilities for masking sensitive field values
// in drift results before they are written to output or logs.
package redact

import "strings"

// DefaultSensitiveKeys contains field names that are redacted by default.
var DefaultSensitiveKeys = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"apikey",
	"credential",
	"private_key",
}

const redactedPlaceholder = "[REDACTED]"

// Options controls redaction behaviour.
type Options struct {
	// ExtraKeys are additional field name substrings to treat as sensitive.
	ExtraKeys []string
	// Disable turns off all redaction when true.
	Disable bool
}

// Result mirrors the minimal drift result shape needed for redaction.
type Result struct {
	Service  string
	Field    string
	Expected string
	Actual   string
	Drifted  bool
}

// Apply returns a copy of results with sensitive field values masked.
func Apply(results []Result, opts Options) []Result {
	if opts.Disable {
		return results
	}

	keys := append(DefaultSensitiveKeys, opts.ExtraKeys...)
	out := make([]Result, len(results))
	for i, r := range results {
		if isSensitive(r.Field, keys) {
			r.Expected = redactedPlaceholder
			r.Actual = redactedPlaceholder
		}
		out[i] = r
	}
	return out
}

// isSensitive reports whether the field name contains any sensitive key substring.
func isSensitive(field string, keys []string) bool {
	lower := strings.ToLower(field)
	for _, k := range keys {
		if strings.Contains(lower, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

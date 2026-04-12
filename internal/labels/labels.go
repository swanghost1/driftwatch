// Package labels provides utilities for matching and filtering drift results
// by arbitrary key=value label pairs attached to declared services.
package labels

import (
	"fmt"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// ParseLabel parses a "key=value" string into its components.
// It returns an error if the format is invalid.
func ParseLabel(raw string) (key, value string, err error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", "", fmt.Errorf("invalid label %q: expected key=value", raw)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

// Match reports whether the result's labels contain all of the required pairs.
// Matching is case-insensitive on both key and value.
func Match(result drift.Result, required map[string]string) bool {
	for k, v := range required {
		got, ok := findLabel(result.Labels, k)
		if !ok || !strings.EqualFold(got, v) {
			return false
		}
	}
	return true
}

// Filter returns only the results whose labels satisfy all required pairs.
func Filter(results []drift.Result, required map[string]string) []drift.Result {
	if len(required) == 0 {
		return results
	}
	out := make([]drift.Result, 0, len(results))
	for _, r := range results {
		if Match(r, required) {
			out = append(out, r)
		}
	}
	return out
}

// ParseAll parses a slice of "key=value" strings into a map.
// The first parse error encountered is returned.
func ParseAll(raws []string) (map[string]string, error) {
	m := make(map[string]string, len(raws))
	for _, raw := range raws {
		k, v, err := ParseLabel(raw)
		if err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, nil
}

func findLabel(labels map[string]string, key string) (string, bool) {
	for k, v := range labels {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return "", false
}

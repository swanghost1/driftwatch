// Package diff provides utilities for computing field-level differences
// between declared and live service state.
package diff

import (
	"fmt"
	"strings"
)

// Field represents a single field-level difference between declared and live state.
type Field struct {
	Name     string
	Declared string
	Live     string
}

// String returns a human-readable representation of the field diff.
func (f Field) String() string {
	return fmt.Sprintf("%s: declared=%q live=%q", f.Name, f.Declared, f.Live)
}

// Result holds all field differences for a single service.
type Result struct {
	Service string
	Fields  []Field
}

// HasDrift reports whether the result contains any field differences.
func (r Result) HasDrift() bool {
	return len(r.Fields) > 0
}

// Summary returns a compact, single-line summary of all drifted fields.
func (r Result) Summary() string {
	if !r.HasDrift() {
		return fmt.Sprintf("%s: no drift", r.Service)
	}
	parts := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		parts[i] = f.String()
	}
	return fmt.Sprintf("%s: %s", r.Service, strings.Join(parts, "; "))
}

// Compare returns the field-level differences between two arbitrary string maps
// representing declared and live state for a named service.
func Compare(service string, declared, live map[string]string) Result {
	res := Result{Service: service}

	for key, declaredVal := range declared {
		liveVal, ok := live[key]
		if !ok || liveVal != declaredVal {
			actual := liveVal
			if !ok {
				actual = "<missing>"
			}
			res.Fields = append(res.Fields, Field{
				Name:     key,
				Declared: declaredVal,
				Live:     actual,
			})
		}
	}

	for key, liveVal := range live {
		if _, ok := declared[key]; !ok {
			res.Fields = append(res.Fields, Field{
				Name:     key,
				Declared: "<missing>",
				Live:     liveVal,
			})
		}
	}

	return res
}

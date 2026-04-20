// Package cascade propagates drift results through a chain of processors,
// applying each stage in order and passing the output of one as the input
// to the next.
package cascade

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Stage is a single step in a processing pipeline.
type Stage struct {
	// Name identifies the stage for reporting purposes.
	Name string
	// Fn applies the stage transformation to a slice of results.
	Fn func([]drift.Result) ([]drift.Result, error)
}

// Options controls cascade behaviour.
type Options struct {
	// StopOnError halts the pipeline if any stage returns an error.
	StopOnError bool
	// Verbose writes per-stage counts to w when set.
	Verbose bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		StopOnError: true,
		Verbose:     false,
	}
}

// Apply runs results through each stage in order, returning the final
// slice and a list of any non-fatal errors encountered when StopOnError
// is false.
func Apply(results []drift.Result, stages []Stage, opts Options, w io.Writer) ([]drift.Result, []error) {
	current := results
	var errs []error

	for _, s := range stages {
		out, err := s.Fn(current)
		if err != nil {
			if opts.StopOnError {
				return current, append(errs, fmt.Errorf("stage %q: %w", s.Name, err))
			}
			errs = append(errs, fmt.Errorf("stage %q: %w", s.Name, err))
			continue
		}
		if opts.Verbose && w != nil {
			fmt.Fprintf(w, "  %-20s %d → %d results\n", s.Name, len(current), len(out))
		}
		current = out
	}

	return current, errs
}

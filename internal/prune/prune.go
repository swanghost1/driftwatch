// Package prune removes drift results that match specified criteria,
// such as results older than a given age or belonging to excluded services.
package prune

import (
	"io"
	"strings"
	"time"

	"fmt"
)

// Options controls which results are pruned.
type Options struct {
	// OlderThan removes results whose DetectedAt is before now minus this duration.
	OlderThan time.Duration
	// Services is a list of service name substrings; matching results are removed.
	Services []string
	// OnlyClean removes results that have no drift.
	OnlyClean bool
}

// Result is a minimal interface for a drift check result.
type Result struct {
	Service    string
	Field      string
	Drifted    bool
	DetectedAt time.Time
}

// Apply removes results from the slice according to opts.
// Results that do not match any prune criterion are retained.
func Apply(results []Result, opts Options) []Result {
	if len(results) == 0 {
		return results
	}

	cutoff := time.Time{}
	if opts.OlderThan > 0 {
		cutoff = time.Now().Add(-opts.OlderThan)
	}

	kept := results[:0]
	for _, r := range results {
		if shouldPrune(r, opts, cutoff) {
			continue
		}
		kept = append(kept, r)
	}
	return kept
}

func shouldPrune(r Result, opts Options, cutoff time.Time) bool {
	if opts.OnlyClean && !r.Drifted {
		return true
	}
	if !cutoff.IsZero() && !r.DetectedAt.IsZero() && r.DetectedAt.Before(cutoff) {
		return true
	}
	for _, svc := range opts.Services {
		if strings.Contains(strings.ToLower(r.Service), strings.ToLower(svc)) {
			return true
		}
	}
	return false
}

// Write outputs a summary of the prune operation to w.
func Write(w io.Writer, before, after []Result) {
	removed := len(before) - len(after)
	fmt.Fprintf(w, "prune: %d results before, %d after, %d removed\n", len(before), len(after), removed)
}

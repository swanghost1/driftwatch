// Package batch provides utilities for processing drift results in
// fixed-size chunks, useful when forwarding results to external systems
// that impose payload size limits.
package batch

import (
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

// Options controls how batching is applied.
type Options struct {
	// Size is the maximum number of results per batch. Zero means no batching.
	Size int
}

// DefaultOptions returns a sensible default batch size.
func DefaultOptions() Options {
	return Options{Size: 50}
}

// Batch is a slice of drift results representing a single chunk.
type Batch []drift.Result

// Apply splits results into batches of at most opts.Size entries.
// If opts.Size is zero or negative, a single batch containing all
// results is returned.
func Apply(results []drift.Result, opts Options) []Batch {
	if len(results) == 0 {
		return nil
	}
	if opts.Size <= 0 {
		return []Batch{results}
	}
	var batches []Batch
	for i := 0; i < len(results); i += opts.Size {
		end := i + opts.Size
		if end > len(results) {
			end = len(results)
		}
		batches = append(batches, Batch(results[i:end]))
	}
	return batches
}

// Write prints a human-readable summary of the batches to w.
func Write(w io.Writer, batches []Batch) {
	fmt.Fprintf(w, "batches: %d\n", len(batches))
	for i, b := range batches {
		drifted := 0
		for _, r := range b {
			if r.Drifted {
				drifted++
			}
		}
		fmt.Fprintf(w, "  [%d] total=%-4d drifted=%d\n", i+1, len(b), drifted)
	}
}

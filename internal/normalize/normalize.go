// Package normalize provides utilities for standardising drift result
// field values before comparison or display, ensuring consistent casing,
// whitespace trimming, and image tag canonicalisation.
package normalize

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Options controls which normalization steps are applied.
type Options struct {
	// TrimWhitespace strips leading/trailing whitespace from string fields.
	TrimWhitespace bool
	// LowercaseImage canonicalises image references to lowercase.
	LowercaseImage bool
	// CanonicaliseTag expands a missing tag to ":latest".
	CanonicaliseTag bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		TrimWhitespace:  true,
		LowercaseImage:  true,
		CanonicaliseTag: true,
	}
}

// Apply runs the configured normalization steps over a slice of drift results,
// returning a new slice with normalised values. The originals are not mutated.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, len(results))
	for i, r := range results {
		out[i] = normaliseResult(r, opts)
	}
	return out
}

func normaliseResult(r drift.Result, opts Options) drift.Result {
	if opts.TrimWhitespace {
		r.Service = strings.TrimSpace(r.Service)
		r.Field = strings.TrimSpace(r.Field)
		r.Expected = strings.TrimSpace(r.Expected)
		r.Actual = strings.TrimSpace(r.Actual)
	}
	if r.Field == "image" {
		if opts.LowercaseImage {
			r.Expected = strings.ToLower(r.Expected)
			r.Actual = strings.ToLower(r.Actual)
		}
		if opts.CanonicaliseTag {
			r.Expected = canonicaliseImage(r.Expected)
			r.Actual = canonicaliseImage(r.Actual)
		}
	}
	return r
}

// canonicaliseImage appends ":latest" when no tag or digest is present.
func canonicaliseImage(image string) string {
	if image == "" {
		return image
	}
	// A colon after the last slash indicates a tag or port; @ indicates a digest.
	base := image
	if idx := strings.LastIndex(image, "/"); idx >= 0 {
		base = image[idx+1:]
	}
	if !strings.Contains(base, ":") && !strings.Contains(image, "@") {
		return image + ":latest"
	}
	return image
}

// Package annotate attaches free-form key/value annotations to drift results,
// allowing downstream consumers to carry extra metadata (e.g. owner, team,
// ticket) without modifying the core Result type.
package annotate

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Options controls which annotations are applied and to which services.
type Options struct {
	// Global annotations are added to every result.
	Global map[string]string

	// PerService maps a service name (case-insensitive substring match) to a
	// set of annotations that are merged on top of the global set.
	PerService map[string]map[string]string
}

// Apply returns a copy of results with annotations merged into each entry's
// Annotations map. Existing annotations are never overwritten.
func Apply(results []drift.Result, opts Options) []drift.Result {
	out := make([]drift.Result, len(results))
	for i, r := range results {
		if r.Annotations == nil {
			r.Annotations = make(map[string]string)
		}
		merge(r.Annotations, opts.Global, false)
		for svc, ann := range opts.PerService {
			if strings.Contains(strings.ToLower(r.Service), strings.ToLower(svc)) {
				merge(r.Annotations, ann, false)
			}
		}
		out[i] = r
	}
	return out
}

// merge copies src key/value pairs into dst. If overwrite is false, existing
// keys in dst are left unchanged.
func merge(dst, src map[string]string, overwrite bool) {
	for k, v := range src {
		if _, exists := dst[k]; !exists || overwrite {
			dst[k] = v
		}
	}
}

// Package annotate provides helpers for attaching arbitrary key/value
// annotations to drift detection results.
//
// Annotations are stored in the Result.Annotations map and are propagated
// through the rest of the pipeline unchanged. They can be used to carry
// organisational metadata such as team ownership, ticket references or
// environment labels.
//
// Usage:
//
//	opts := annotate.Options{
//	    Global: map[string]string{"env": "production"},
//	    PerService: map[string]map[string]string{
//	        "payments": {"owner": "payments-team", "ticket": "OPS-42"},
//	    },
//	}
//	annotated := annotate.Apply(results, opts)
package annotate

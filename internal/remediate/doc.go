// Package remediate analyses drift detection results and produces actionable
// remediation suggestions.
//
// For each drifted field, Suggest returns a Suggestion that includes a
// human-readable Action string modelled on common kubectl commands, making it
// straightforward to copy-paste corrections directly into a terminal.
//
// Supported fields with tailored suggestions:
//
//	"image"    → kubectl set image
//	"replicas" → kubectl scale
//	"env:<KEY>" → kubectl set env
//
// All other fields fall back to a freeform manual-reconcile comment so that
// no drift goes unacknowledged.
//
// Usage:
//
//	suggestions := remediate.Suggest(results)
//	for _, s := range suggestions {
//		fmt.Println(s.Action)
//	}
package remediate

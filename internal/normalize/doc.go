// Package normalize standardises drift result field values before they are
// compared, stored, or displayed.
//
// Normalization is a pre-processing step that removes superficial differences
// between declared and live state, such as inconsistent casing in image
// references or missing ":latest" tags, so that the drift detector can focus
// on semantically meaningful divergence.
//
// # Usage
//
//	opts := normalize.DefaultOptions()
//	normalised := normalize.Apply(results, opts)
//
// DefaultOptions enables all normalization steps:
//   - TrimWhitespace — strips leading/trailing whitespace from every string field.
//   - LowercaseImage — converts image Expected/Actual values to lowercase.
//   - CanonicaliseTag — appends ":latest" to image references that carry no tag
//     or digest, matching the implicit behaviour of container runtimes.
//
// Apply never mutates the input slice; it always returns a new slice.
package normalize

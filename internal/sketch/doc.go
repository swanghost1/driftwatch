// Package sketch condenses drift results into a per-service summary,
// providing a high-level overview of which services are drifted and
// which fields are responsible.
//
// Use Compute to build a []Entry from drift results, then Write or
// WriteJSON to render the summary to any io.Writer.
//
// Example:
//
//	entries := sketch.Compute(results)
//	sketch.Write(os.Stdout, entries)
package sketch

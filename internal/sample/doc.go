// Package sample implements probabilistic sampling of drift run results.
//
// When driftwatch runs frequently, persisting every result can produce
// large amounts of data. The sample package provides two complementary
// mechanisms to manage this:
//
//   - ShouldRecord decides, before writing anything, whether the current
//     run should be persisted at all. A configurable Rate controls the
//     probability, while AlwaysSampleDrifted guarantees that runs
//     containing at least one drifted service are always retained.
//
//   - Apply trims an already-collected result slice to at most n entries,
//     preferring drifted results over clean ones so that signal is not
//     discarded in favour of noise.
//
// Typical usage:
//
//	opts := sample.DefaultOptions()
//	if sample.ShouldRecord(results, opts) {
//		sampled := sample.Apply(results, 50)
//		// persist sampled ...
//	}
package sample

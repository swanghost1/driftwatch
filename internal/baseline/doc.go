// Package baseline provides functionality for saving, loading, and comparing
// drift detection baselines.
//
// A baseline captures the set of drift results at a known point in time and
// can be used in subsequent runs to surface only newly introduced drift,
// reducing noise from pre-existing, acknowledged discrepancies.
//
// Typical usage:
//
//	// Save current drift results as the new baseline
//	err := baseline.Save(".driftwatch/baseline.json", "release-1.4.2", results)
//
//	// On the next run, load the baseline and compare
//	entry, err := baseline.Load(".driftwatch/baseline.json")
//	novel := baseline.Compare(entry, newResults)
//	if len(novel) > 0 {
//		// report only newly introduced drift
//	}
package baseline

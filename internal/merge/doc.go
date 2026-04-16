// Package merge provides utilities for combining multiple drift result
// slices produced by independent detection runs or data sources.
//
// # Overview
//
// When driftwatch is run against several configuration sources — for
// example a local YAML file and a remote registry snapshot — the
// individual result sets may overlap. The merge package resolves those
// overlaps using a simple precedence rule: a drifted result always
// takes priority over a clean result for the same service+field pair.
//
// # Usage
//
//	merged := merge.Apply(resultsA, resultsB)
//	merge.Write(os.Stdout, merged)
//
// The returned slice is deduplicated and sorted by service name then
// field name, making the output deterministic across runs.
package merge

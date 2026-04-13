// Package compare provides field-level comparison utilities for drift results,
// allowing callers to compute symmetric differences between two sets of drift
// results and identify newly resolved or newly introduced drifts.
package compare

import (
	"fmt"

	"github.com/example/driftwatch/internal/drift"
)

// ChangeKind describes how a drift result changed between two snapshots.
type ChangeKind string

const (
	ChangeIntroduced ChangeKind = "introduced"
	ChangeResolved   ChangeKind = "resolved"
	ChangeUnchanged  ChangeKind = "unchanged"
)

// Change represents a single drift result and how it changed.
type Change struct {
	Result drift.Result
	Kind   ChangeKind
}

// Diff computes the symmetric difference between a previous and current set of
// drift results. Results are keyed by "service:field" to detect transitions.
func Diff(previous, current []drift.Result) []Change {
	prev := index(previous)
	curr := index(current)

	var changes []Change

	for key, r := range curr {
		if _, existed := prev[key]; existed {
			changes = append(changes, Change{Result: r, Kind: ChangeUnchanged})
		} else {
			changes = append(changes, Change{Result: r, Kind: ChangeIntroduced})
		}
	}

	for key, r := range prev {
		if _, exists := curr[key]; !exists {
			changes = append(changes, Change{Result: r, Kind: ChangeResolved})
		}
	}

	return changes
}

// OnlyIntroduced filters a slice of Changes to those with kind ChangeIntroduced.
func OnlyIntroduced(changes []Change) []Change {
	return filterKind(changes, ChangeIntroduced)
}

// OnlyResolved filters a slice of Changes to those with kind ChangeResolved.
func OnlyResolved(changes []Change) []Change {
	return filterKind(changes, ChangeResolved)
}

func filterKind(changes []Change, kind ChangeKind) []Change {
	var out []Change
	for _, c := range changes {
		if c.Kind == kind {
			out = append(out, c)
		}
	}
	return out
}

func index(results []drift.Result) map[string]drift.Result {
	m := make(map[string]drift.Result, len(results))
	for _, r := range results {
		if r.Drifted {
			key := fmt.Sprintf("%s:%s", r.Service, r.Field)
			m[key] = r
		}
	}
	return m
}

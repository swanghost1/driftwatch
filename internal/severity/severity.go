// Package severity classifies drift results by severity level based on
// which fields have drifted and optional per-field weight configuration.
package severity

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/example/driftwatch/internal/drift"
)

// Level represents a drift severity classification.
type Level string

const (
	LevelNone     Level = "none"
	LevelLow      Level = "low"
	LevelMedium   Level = "medium"
	LevelHigh     Level = "high"
	LevelCritical Level = "critical"
)

// Config maps field names to severity levels. Fields not listed default to Low.
type Config map[string]Level

// DefaultConfig returns a sensible default field-to-severity mapping.
func DefaultConfig() Config {
	return Config{
		"image":    LevelCritical,
		"replicas": LevelHigh,
		"env":      LevelMedium,
		"port":     LevelLow,
	}
}

// Result pairs a drift result with its computed severity.
type Result struct {
	Drift    drift.Result
	Severity Level
}

// Classify assigns a severity level to each drift result.
// Results with no drift are assigned LevelNone.
func Classify(results []drift.Result, cfg Config) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Drift:    r,
			Severity: classify(r, cfg),
		})
	}
	return out
}

func classify(r drift.Result, cfg Config) Level {
	if !r.Drifted {
		return LevelNone
	}
	highest := LevelLow
	for _, d := range r.Diffs {
		lvl, ok := cfg[d.Field]
		if !ok {
			lvl = LevelLow
		}
		if levelRank(lvl) > levelRank(highest) {
			highest = lvl
		}
	}
	return highest
}

func levelRank(l Level) int {
	switch l {
	case LevelNone:
		return 0
	case LevelLow:
		return 1
	case LevelMedium:
		return 2
	case LevelHigh:
		return 3
	case LevelCritical:
		return 4
	default:
		return 1
	}
}

// Write renders a severity-annotated table to w.
func Write(w io.Writer, results []Result) error {
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return levelRank(sorted[i].Severity) > levelRank(sorted[j].Severity)
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSEVERITY\tDRIFTED")
	for _, r := range sorted {
		drifted := "no"
		if r.Drift.Drifted {
			drifted = "yes"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.Drift.Service, r.Severity, drifted)
	}
	return tw.Flush()
}

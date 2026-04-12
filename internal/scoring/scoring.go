// Package scoring computes a numeric drift health score for a set of
// detection results. A score of 100 means no drift; each drifted service
// reduces the score proportionally.
package scoring

import (
	"fmt"
	"io"
	"math"

	"github.com/example/driftwatch/internal/drift"
)

// Score holds the computed health score and supporting statistics.
type Score struct {
	Total   int     `json:"total"`
	Drifted int     `json:"drifted"`
	Clean   int     `json:"clean"`
	Value   float64 `json:"score"` // 0‥100
	Grade   string  `json:"grade"`
}

// Compute calculates the drift health score from a slice of results.
func Compute(results []drift.Result) Score {
	total := len(results)
	if total == 0 {
		return Score{Grade: "A", Value: 100}
	}

	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	clean := total - drifted
	value := math.Round((float64(clean)/float64(total))*10000) / 100

	return Score{
		Total:   total,
		Drifted: drifted,
		Clean:   clean,
		Value:   value,
		Grade:   grade(value),
	}
}

// grade converts a numeric score to a letter grade.
func grade(score float64) string {
	switch {
	case score >= 95:
		return "A"
	case score >= 80:
		return "B"
	case score >= 65:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

// Write renders a human-readable score summary to w.
func Write(w io.Writer, s Score) {
	fmt.Fprintf(w, "Drift Health Score: %.2f / 100  (Grade: %s)\n", s.Value, s.Grade)
	fmt.Fprintf(w, "  Services checked : %d\n", s.Total)
	fmt.Fprintf(w, "  Clean            : %d\n", s.Clean)
	fmt.Fprintf(w, "  Drifted          : %d\n", s.Drifted)
}

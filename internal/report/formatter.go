package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Summary holds aggregated drift report data.
type Summary struct {
	TotalServices int
	DriftedServices int
	Results []drift.Result
}

// NewSummary builds a Summary from a slice of drift results.
func NewSummary(results []drift.Result) Summary {
	drifted := 0
	for _, r := range results {
		if r.HasDrift() {
			drifted++
		}
	}
	return Summary{
		TotalServices:   len(results),
		DriftedServices: drifted,
		Results:         results,
	}
}

// Write renders the summary to w in the requested format.
func Write(w io.Writer, s Summary, f Format) error {
	switch f {
	case FormatJSON:
		return writeJSON(w, s)
	default:
		return writeText(w, s)
	}
}

func writeText(w io.Writer, s Summary) error {
	fmt.Fprintf(w, "Drift Report\n%s\n", strings.Repeat("=", 40))
	fmt.Fprintf(w, "Services checked : %d\n", s.TotalServices)
	fmt.Fprintf(w, "Services drifted : %d\n\n", s.DriftedServices)

	for _, r := range s.Results {
		if !r.HasDrift() {
			fmt.Fprintf(w, "[OK]    %s\n", r.ServiceName)
			continue
		}
		fmt.Fprintf(w, "[DRIFT] %s\n", r.ServiceName)
		for _, d := range r.Diffs {
			fmt.Fprintf(w, "        - %s: expected %q, got %q\n", d.Field, d.Expected, d.Actual)
		}
	}
	return nil
}

func writeJSON(w io.Writer, s Summary) error {
	// Manual JSON to avoid importing encoding/json for a lightweight formatter.
	fmt.Fprintf(w, "{\n  \"total_services\": %d,\n  \"drifted_services\": %d,\n  \"results\": [\n", s.TotalServices, s.DriftedServices)
	for i, r := range s.Results {
		comma := ","
		if i == len(s.Results)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "    {\"service\": %q, \"drifted\": %v, \"diffs\": %d}%s\n",
			r.ServiceName, r.HasDrift(), len(r.Diffs), comma)
	}
	fmt.Fprintf(w, "  ]\n}\n")
	return nil
}

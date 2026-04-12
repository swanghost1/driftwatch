// Package export provides functionality for exporting drift results
// to external formats such as CSV and Markdown tables.
package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Format represents an export output format.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Write serialises results in the requested format to w.
func Write(w io.Writer, results []drift.Result, format Format) error {
	switch format {
	case FormatCSV:
		return writeCSV(w, results)
	case FormatMarkdown:
		return writeMarkdown(w, results)
	default:
		return fmt.Errorf("export: unknown format %q", format)
	}
}

func writeCSV(w io.Writer, results []drift.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"service", "status", "field", "declared", "live"}); err != nil {
		return err
	}
	for _, r := range results {
		if len(r.Drifts) == 0 {
			if err := cw.Write([]string{r.Service, "ok", "", "", ""}); err != nil {
				return err
			}
			continue
		}
		for _, d := range r.Drifts {
			if err := cw.Write([]string{r.Service, "drift", d.Field, d.Declared, d.Live}); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(w io.Writer, results []drift.Result) error {
	fmt.Fprintln(w, "| Service | Status | Field | Declared | Live |")
	fmt.Fprintln(w, "|---------|--------|-------|----------|------|")
	for _, r := range results {
		if len(r.Drifts) == 0 {
			fmt.Fprintf(w, "| %s | ok | | | |\n", r.Service)
			continue
		}
		for _, d := range r.Drifts {
			fmt.Fprintf(w, "| %s | drift | %s | %s | %s |\n",
				r.Service,
				escapeMD(d.Field),
				escapeMD(d.Declared),
				escapeMD(d.Live),
			)
		}
	}
	return nil
}

func escapeMD(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}

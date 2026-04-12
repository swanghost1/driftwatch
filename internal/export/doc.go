// Package export serialises drift results to portable external formats.
//
// Supported formats:
//
//	"csv"      — comma-separated values, suitable for spreadsheet import.
//	"markdown" — GitHub-flavoured Markdown table, suitable for reports.
//
// Usage:
//
//	var buf bytes.Buffer
//	if err := export.Write(&buf, results, export.FormatCSV); err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Print(buf.String())
package export

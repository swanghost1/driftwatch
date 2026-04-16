// Package pivot provides cross-tabulation of drift results along a chosen
// axis — either "service" or "field".
//
// Use Compute to build a Table from a slice of drift.Result values, then
// Write or WriteJSON to render it.
//
// Example:
//
//	tbl := pivot.Compute(results, pivot.AxisService)
//	pivot.Write(os.Stdout, tbl)
//
// The resulting table lists each key (service name or field name) together
// with the total number of checks, how many were drifted, how many were
// clean, and the drift percentage. Rows are sorted by descending drift
// count so the most problematic entries appear first.
package pivot

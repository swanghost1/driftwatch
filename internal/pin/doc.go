// Package pin implements result pinning for driftwatch.
//
// A pinned result is a drift finding that has been explicitly acknowledged by
// an operator. Pinning suppresses the drift from reports and counts without
// altering the underlying declared or live state — it simply marks the
// discrepancy as intentional or temporarily acceptable.
//
// Pin rules are stored as a JSON file and support:
//
//   - Service-scoped rules (exact match or wildcard "*")
//   - Optional field scoping (empty field pins all fields for the service)
//   - Optional expiry timestamps — expired rules are ignored automatically
//   - A free-text reason field for audit purposes
//
// Typical usage:
//
//	store, err := pin.LoadStore(".driftwatch/pins.json")
//	results = pin.Apply(results, store)
//
// Use WriteText or WriteJSON to display currently active pin rules.
package pin

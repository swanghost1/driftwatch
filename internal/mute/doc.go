// Package mute provides time-bounded silencing of drift results.
//
// A mute rule suppresses drift notifications and report entries for a
// specific service (and optionally a specific field) until a declared
// expiry time. This is useful during planned maintenance windows or
// known-good rollouts where transient drift is expected.
//
// Rules are persisted as a JSON store on disk so they survive across
// driftwatch invocations. Expired rules are ignored automatically; a
// separate housekeeping step can prune them via SaveStore.
//
// Example usage:
//
//	store, _ := mute.LoadStore(".driftwatch/mutes.json")
//	filtered := mute.Apply(results, store, time.Now())
package mute

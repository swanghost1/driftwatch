// Package rank provides priority-based ordering of drift results.
//
// Results can be sorted by one or more criteria applied left-to-right:
//
//	- drifted_first  — drifted results appear before clean ones
//	- severity       — critical → high → medium → low → unknown
//	- service        — alphabetical by service name
//	- field          — alphabetical by field name
//
// Example usage:
//
//	opts := rank.Options{
//		Criteria: []rank.Criterion{rank.ByDriftOnly, rank.BySeverity},
//	}
//	ordered := rank.Apply(results, opts)
//	rank.Write(os.Stdout, ordered)
//
// Apply never modifies the input slice; it returns a new sorted copy.
package rank

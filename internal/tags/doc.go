// Package tags provides grouping and filtering utilities for drift results
// based on arbitrary key-value tags attached to services in the driftwatch
// configuration.
//
// Tags are declared per-service in the YAML config:
//
//	services:
//	  - name: payments-api
//	    tags:
//	      env: prod
//	      team: platform
//
// GroupByTag partitions a slice of drift.Result values by the value of a
// chosen tag key, enabling per-environment or per-team reporting.
//
// FilterByTag narrows results to those whose chosen tag matches a given
// value, complementing the broader service-name filter in internal/filter.
package tags

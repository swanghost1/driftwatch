// Package alert evaluates drift results against user-defined threshold rules
// and produces structured alerts when those thresholds are exceeded.
//
// # Rules
//
// A Rule specifies a minimum number of drifted services (MinDrifted) and a
// severity Level (warn or error). When the number of drifted services in a
// result set meets or exceeds MinDrifted, an Alert is produced.
//
// Multiple rules may be defined; each is evaluated independently so that
// different severities can be triggered by different thresholds.
//
// # Usage
//
//	alerts := alert.Evaluate(results, rules)
//	alert.Write(os.Stdout, alerts)
//
// Alerts can be forwarded to notification backends (e.g. internal/notify)
// after evaluation.
package alert

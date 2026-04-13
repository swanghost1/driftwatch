package redact_test

import (
	"testing"

	"github.com/driftwatch/internal/redact"
)

func makeResults() []redact.Result {
	return []redact.Result{
		{Service: "api", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24", Drifted: true},
		{Service: "api", Field: "env.PASSWORD", Expected: "hunter2", Actual: "letmein", Drifted: true},
		{Service: "api", Field: "env.API_KEY", Expected: "abc123", Actual: "xyz789", Drifted: true},
		{Service: "api", Field: "replicas", Expected: "3", Actual: "2", Drifted: true},
	}
}

func TestApply_NonSensitiveField_NotRedacted(t *testing.T) {
	results := makeResults()
	out := redact.Apply(results, redact.Options{})
	for _, r := range out {
		if r.Field == "image" {
			if r.Expected != "nginx:1.25" || r.Actual != "nginx:1.24" {
				t.Errorf("image field should not be redacted")
			}
		}
	}
}

func TestApply_SensitiveField_Redacted(t *testing.T) {
	out := redact.Apply(makeResults(), redact.Options{})
	for _, r := range out {
		if r.Field == "env.PASSWORD" || r.Field == "env.API_KEY" {
			if r.Expected != "[REDACTED]" || r.Actual != "[REDACTED]" {
				t.Errorf("field %q should be redacted, got expected=%q actual=%q", r.Field, r.Expected, r.Actual)
			}
		}
	}
}

func TestApply_Disable_SkipsRedaction(t *testing.T) {
	out := redact.Apply(makeResults(), redact.Options{Disable: true})
	for _, r := range out {
		if r.Field == "env.PASSWORD" && r.Expected == "[REDACTED]" {
			t.Error("redaction should be disabled")
		}
	}
}

func TestApply_ExtraKeys_RedactsCustomField(t *testing.T) {
	results := []redact.Result{
		{Service: "svc", Field: "env.DB_PASS", Expected: "secret", Actual: "other", Drifted: true},
	}
	out := redact.Apply(results, redact.Options{ExtraKeys: []string{"db_pass"}})
	if out[0].Expected != "[REDACTED]" {
		t.Errorf("expected DB_PASS to be redacted")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	results := []redact.Result{
		{Service: "svc", Field: "env.Auth_Token", Expected: "tok", Actual: "tok2", Drifted: true},
	}
	out := redact.Apply(results, redact.Options{})
	if out[0].Expected != "[REDACTED]" {
		t.Errorf("expected Auth_Token to be redacted (case-insensitive)")
	}
}

func TestApply_PreservesNonSensitiveValues(t *testing.T) {
	results := []redact.Result{
		{Service: "svc", Field: "replicas", Expected: "3", Actual: "2", Drifted: true},
	}
	out := redact.Apply(results, redact.Options{})
	if out[0].Expected != "3" || out[0].Actual != "2" {
		t.Errorf("replicas should not be redacted")
	}
}

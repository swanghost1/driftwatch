package mask_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/mask"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "svc-a", Field: "image", Declared: "nginx:1.25", Live: "nginx:1.24", Drifted: true},
		{Service: "svc-a", Field: "password", Declared: "hunter2", Live: "hunter2", Drifted: false},
		{Service: "svc-b", Field: "TOKEN", Declared: "abc123", Live: "xyz789", Drifted: true},
		{Service: "svc-b", Field: "replicas", Declared: "3", Live: "2", Drifted: true},
	}
}

func TestApply_NonSensitiveField_NotMasked(t *testing.T) {
	results := makeResults()
	out := mask.Apply(results, mask.DefaultOptions())

	for _, r := range out {
		if r.Field == "image" {
			if r.Live == "***" || r.Declared == "***" {
				t.Errorf("image field should not be masked")
			}
		}
	}
}

func TestApply_SensitiveField_MaskedBothValues(t *testing.T) {
	results := makeResults()
	out := mask.Apply(results, mask.DefaultOptions())

	for _, r := range out {
		if r.Field == "password" {
			if r.Live != "***" || r.Declared != "***" {
				t.Errorf("password field values should be masked, got live=%q declared=%q", r.Live, r.Declared)
			}
		}
	}
}

func TestApply_CaseInsensitiveFieldMatch(t *testing.T) {
	results := makeResults()
	out := mask.Apply(results, mask.DefaultOptions())

	for _, r := range out {
		if r.Field == "TOKEN" {
			if r.Live != "***" || r.Declared != "***" {
				t.Errorf("TOKEN field should be masked case-insensitively")
			}
		}
	}
}

func TestApply_Disabled_SkipsMasking(t *testing.T) {
	results := makeResults()
	opts := mask.Options{Fields: []string{"password", "token"}, Disabled: true}
	out := mask.Apply(results, opts)

	for i, r := range out {
		if r.Live != results[i].Live || r.Declared != results[i].Declared {
			t.Errorf("disabled masking should leave values unchanged for field %q", r.Field)
		}
	}
}

func TestApply_CustomField_IsMasked(t *testing.T) {
	results := []drift.Result{
		{Service: "svc", Field: "db_pass", Declared: "secret", Live: "other", Drifted: true},
	}
	opts := mask.Options{Fields: []string{"db_pass"}}
	out := mask.Apply(results, opts)

	if out[0].Live != "***" || out[0].Declared != "***" {
		t.Errorf("custom field db_pass should be masked")
	}
}

func TestApply_EmptyFields_ReturnsUnchanged(t *testing.T) {
	results := makeResults()
	opts := mask.Options{Fields: []string{}}
	out := mask.Apply(results, opts)

	for i, r := range out {
		if r.Live != results[i].Live {
			t.Errorf("empty field list should not mask anything")
		}
	}
}

func TestApply_PreservesOtherFields(t *testing.T) {
	results := makeResults()
	out := mask.Apply(results, mask.DefaultOptions())

	for i, r := range out {
		if r.Service != results[i].Service || r.Field != results[i].Field || r.Drifted != results[i].Drifted {
			t.Errorf("Apply should not alter Service, Field, or Drifted metadata")
		}
	}
}

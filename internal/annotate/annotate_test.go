package annotate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/annotate"
	"github.com/example/driftwatch/internal/drift"
)

func makeResult(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestApply_GlobalAnnotations_AddedToAll(t *testing.T) {
	results := []drift.Result{makeResult("alpha"), makeResult("beta")}
	opts := annotate.Options{
		Global: map[string]string{"env": "staging"},
	}
	out := annotate.Apply(results, opts)
	for _, r := range out {
		if got := r.Annotations["env"]; got != "staging" {
			t.Errorf("service %s: want env=staging, got %q", r.Service, got)
		}
	}
}

func TestApply_PerServiceAnnotation_OnlyMatchingService(t *testing.T) {
	results := []drift.Result{makeResult("payments-api"), makeResult("inventory")}
	opts := annotate.Options{
		PerService: map[string]map[string]string{
			"payments": {"owner": "payments-team"},
		},
	}
	out := annotate.Apply(results, opts)
	if got := out[0].Annotations["owner"]; got != "payments-team" {
		t.Errorf("want payments-team, got %q", got)
	}
	if got := out[1].Annotations["owner"]; got != "" {
		t.Errorf("inventory should have no owner annotation, got %q", got)
	}
}

func TestApply_ExistingAnnotation_NotOverwritten(t *testing.T) {
	r := makeResult("svc")
	r.Annotations = map[string]string{"env": "prod"}
	opts := annotate.Options{
		Global: map[string]string{"env": "staging"},
	}
	out := annotate.Apply([]drift.Result{r}, opts)
	if got := out[0].Annotations["env"]; got != "prod" {
		t.Errorf("existing annotation should not be overwritten, got %q", got)
	}
}

func TestApply_CaseInsensitiveServiceMatch(t *testing.T) {
	results := []drift.Result{makeResult("Payments-API")}
	opts := annotate.Options{
		PerService: map[string]map[string]string{
			"payments": {"team": "finance"},
		},
	}
	out := annotate.Apply(results, opts)
	if got := out[0].Annotations["team"]; got != "finance" {
		t.Errorf("want finance, got %q", got)
	}
}

func TestApply_EmptyOptions_ReturnsUnchanged(t *testing.T) {
	results := []drift.Result{makeResult("svc")}
	out := annotate.Apply(results, annotate.Options{})
	if len(out[0].Annotations) != 0 {
		t.Errorf("expected no annotations, got %v", out[0].Annotations)
	}
}

func TestWriteText_NoAnnotations_ShowsNone(t *testing.T) {
	var buf bytes.Buffer
	annotate.WriteText(&buf, []drift.Result{makeResult("svc")})
	if !strings.Contains(buf.String(), "none") {
		t.Errorf("expected 'none' in output, got: %s", buf.String())
	}
}

func TestWriteText_WithAnnotations_ContainsKeyValue(t *testing.T) {
	r := makeResult("svc")
	r.Annotations = map[string]string{"owner": "team-a"}
	var buf bytes.Buffer
	annotate.WriteText(&buf, []drift.Result{r})
	if !strings.Contains(buf.String(), "owner") || !strings.Contains(buf.String(), "team-a") {
		t.Errorf("expected key and value in output, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	r := makeResult("svc")
	r.Annotations = map[string]string{"ticket": "OPS-1"}
	var buf bytes.Buffer
	if err := annotate.WriteJSON(&buf, []drift.Result{r}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "OPS-1") {
		t.Errorf("expected ticket value in JSON, got: %s", buf.String())
	}
}

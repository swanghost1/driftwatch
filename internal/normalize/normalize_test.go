package normalize_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/normalize"
)

func makeResult(service, field, expected, actual string, drifted bool) drift.Result {
	return drift.Result{
		Service:  service,
		Field:    field,
		Expected: expected,
		Actual:   actual,
		Drifted:  drifted,
	}
}

func TestApply_TrimWhitespace(t *testing.T) {
	input := []drift.Result{
		makeResult("  svc  ", "  image  ", "  nginx:1.0  ", "  nginx:2.0  ", true),
	}
	out := normalize.Apply(input, normalize.Options{TrimWhitespace: true})
	if out[0].Service != "svc" {
		t.Errorf("expected trimmed service, got %q", out[0].Service)
	}
	if out[0].Expected != "nginx:1.0" {
		t.Errorf("expected trimmed Expected, got %q", out[0].Expected)
	}
}

func TestApply_LowercaseImage(t *testing.T) {
	input := []drift.Result{
		makeResult("svc", "image", "NGINX:1.0", "NGINX:2.0", true),
	}
	out := normalize.Apply(input, normalize.Options{LowercaseImage: true})
	if out[0].Expected != "nginx:1.0" {
		t.Errorf("expected lowercase Expected, got %q", out[0].Expected)
	}
	if out[0].Actual != "nginx:2.0" {
		t.Errorf("expected lowercase Actual, got %q", out[0].Actual)
	}
}

func TestApply_CanonicaliseTag_NoTag(t *testing.T) {
	input := []drift.Result{
		makeResult("svc", "image", "nginx", "nginx", false),
	}
	out := normalize.Apply(input, normalize.Options{CanonicaliseTag: true})
	if out[0].Expected != "nginx:latest" {
		t.Errorf("expected :latest appended, got %q", out[0].Expected)
	}
}

func TestApply_CanonicaliseTag_WithTag_NoChange(t *testing.T) {
	input := []drift.Result{
		makeResult("svc", "image", "nginx:1.25", "nginx:1.26", true),
	}
	out := normalize.Apply(input, normalize.Options{CanonicaliseTag: true})
	if out[0].Expected != "nginx:1.25" {
		t.Errorf("expected unchanged tag, got %q", out[0].Expected)
	}
}

func TestApply_CanonicaliseTag_WithDigest_NoChange(t *testing.T) {
	image := "nginx@sha256:abc123"
	input := []drift.Result{
		makeResult("svc", "image", image, image, false),
	}
	out := normalize.Apply(input, normalize.Options{CanonicaliseTag: true})
	if out[0].Expected != image {
		t.Errorf("expected digest image unchanged, got %q", out[0].Expected)
	}
}

func TestApply_NonImageField_NotLowercased(t *testing.T) {
	input := []drift.Result{
		makeResult("svc", "env", "KEY=VALUE", "KEY=OTHER", true),
	}
	out := normalize.Apply(input, normalize.Options{LowercaseImage: true})
	if out[0].Expected != "KEY=VALUE" {
		t.Errorf("non-image field should not be lowercased, got %q", out[0].Expected)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := []drift.Result{
		makeResult("  svc  ", "image", "NGINX", "NGINX", false),
	}
	orig := input[0].Service
	normalize.Apply(input, normalize.DefaultOptions())
	if input[0].Service != orig {
		t.Error("Apply must not mutate the input slice")
	}
}

func TestDefaultOptions_AllEnabled(t *testing.T) {
	opts := normalize.DefaultOptions()
	if !opts.TrimWhitespace || !opts.LowercaseImage || !opts.CanonicaliseTag {
		t.Error("DefaultOptions should enable all normalization steps")
	}
}

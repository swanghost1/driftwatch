package evict_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/evict"
)

var ref = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeResult(service string, drifted bool, age time.Duration) drift.Result {
	return drift.Result{
		Service:    service,
		Drifted:    drifted,
		DetectedAt: ref.Add(-age),
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", true, time.Hour),
		makeResult("svc-b", false, 2*time.Hour),
	}
	opts := evict.Options{Reference: ref}
	out := evict.Apply(results, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApply_MaxAge_RemovesOldResults(t *testing.T) {
	results := []drift.Result{
		makeResult("fresh", true, 1*time.Hour),
		makeResult("stale", true, 100*time.Hour),
	}
	opts := evict.Options{MaxAge: 48 * time.Hour, Reference: ref}
	out := evict.Apply(results, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Service != "fresh" {
		t.Errorf("expected fresh, got %s", out[0].Service)
	}
}

func TestApply_OnlyDrifted_RemovesClean(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", true, time.Hour),
		makeResult("svc-b", false, time.Hour),
		makeResult("svc-c", true, time.Hour),
	}
	opts := evict.Options{OnlyDrifted: true, Reference: ref}
	out := evict.Apply(results, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	for _, r := range out {
		if !r.Drifted {
			t.Errorf("clean result survived eviction: %s", r.Service)
		}
	}
}

func TestApply_BothOptions_AppliedTogether(t *testing.T) {
	results := []drift.Result{
		makeResult("drifted-fresh", true, time.Hour),
		makeResult("drifted-stale", true, 200*time.Hour),
		makeResult("clean-fresh", false, time.Hour),
	}
	opts := evict.Options{MaxAge: 48 * time.Hour, OnlyDrifted: true, Reference: ref}
	out := evict.Apply(results, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Service != "drifted-fresh" {
		t.Errorf("unexpected service: %s", out[0].Service)
	}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	out := evict.Apply(nil, evict.DefaultOptions())
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestApply_ZeroMaxAge_SkipsAgeCheck(t *testing.T) {
	results := []drift.Result{
		makeResult("ancient", true, 9999*time.Hour),
	}
	opts := evict.Options{MaxAge: 0, Reference: ref}
	out := evict.Apply(results, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestWrite_ShowsEvictedAndRetained(t *testing.T) {
	var buf bytes.Buffer
	evict.Write(&buf, 10, 7)
	got := buf.String()
	if got != "evict: 3 result(s) removed, 7 retained\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

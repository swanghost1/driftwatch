package window_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/window"
)

func makeResult(service string, detectedAt time.Time, drifted bool) drift.Result {
	r := drift.Result{
		Service:    service,
		Drifted:    drifted,
		DetectedAt: detectedAt,
	}
	return r
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("a", now.Add(-2*time.Hour), true),
		makeResult("b", now.Add(-1*time.Hour), false),
	}
	got := window.Apply(results, window.Options{})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_Since_ExcludesOlder(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("old", now.Add(-3*time.Hour), true),
		makeResult("new", now.Add(-30*time.Minute), true),
	}
	got := window.Apply(results, window.Options{Since: now.Add(-1 * time.Hour)})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Service != "new" {
		t.Errorf("expected new, got %s", got[0].Service)
	}
}

func TestApply_Until_ExcludesNewer(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("recent", now.Add(-10*time.Minute), true),
		makeResult("old", now.Add(-5*time.Hour), false),
	}
	got := window.Apply(results, window.Options{Until: now.Add(-1 * time.Hour)})
	if len(got) != 1 || got[0].Service != "old" {
		t.Errorf("expected only 'old', got %+v", got)
	}
}

func TestApply_SinceAndUntil_BothBoundsApplied(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("too-old", now.Add(-6*time.Hour), true),
		makeResult("in-range", now.Add(-2*time.Hour), true),
		makeResult("too-new", now.Add(-5*time.Minute), true),
	}
	got := window.Apply(results, window.Options{
		Since: now.Add(-4 * time.Hour),
		Until: now.Add(-1 * time.Hour),
	})
	if len(got) != 1 || got[0].Service != "in-range" {
		t.Errorf("expected only 'in-range', got %+v", got)
	}
}

func TestLast_ReturnsRecentResults(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("stale", now.Add(-10*time.Hour), true),
		makeResult("fresh", now.Add(-30*time.Minute), false),
	}
	got := window.Last(results, 2*time.Hour)
	if len(got) != 1 || got[0].Service != "fresh" {
		t.Errorf("expected only 'fresh', got %+v", got)
	}
}

func TestLast_ZeroDuration_ReturnsAll(t *testing.T) {
	now := time.Now()
	results := []drift.Result{
		makeResult("a", now.Add(-10*time.Hour), true),
		makeResult("b", now.Add(-1*time.Minute), false),
	}
	got := window.Last(results, 0)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

package stale_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/stale"
)

var epoch = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func makeResult(service, field string, drifted bool) drift.Result {
	return drift.Result{
		Service: service,
		Field:   field,
		Drifted: drifted,
	}
}

func TestMark_NoDrift_NeverStale(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", false)}
	firstSeen := map[string]time.Time{}
	out := stale.Mark(results, firstSeen, stale.Options{After: time.Hour, Reference: epoch})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Stale {
		t.Error("non-drifted result should not be stale")
	}
}

func TestMark_RecentDrift_NotStale(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", true)}
	firstSeen := map[string]time.Time{"svc-a:image": epoch.Add(-30 * time.Minute)}
	out := stale.Mark(results, firstSeen, stale.Options{After: time.Hour, Reference: epoch})
	if out[0].Stale {
		t.Error("drift younger than threshold should not be stale")
	}
}

func TestMark_OldDrift_MarkedStale(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", true)}
	firstSeen := map[string]time.Time{"svc-a:image": epoch.Add(-2 * time.Hour)}
	out := stale.Mark(results, firstSeen, stale.Options{After: time.Hour, Reference: epoch})
	if !out[0].Stale {
		t.Error("drift older than threshold should be stale")
	}
	if out[0].Age < time.Hour {
		t.Errorf("expected age >= 1h, got %s", out[0].Age)
	}
}

func TestMark_NewDrift_RecordsFirstSeen(t *testing.T) {
	results := []drift.Result{makeResult("svc-b", "replicas", true)}
	out := stale.Mark(results, map[string]time.Time{}, stale.Options{Reference: epoch})
	if out[0].FirstSeen != epoch {
		t.Errorf("expected FirstSeen=%v, got %v", epoch, out[0].FirstSeen)
	}
}

func TestUpdateFirstSeen_AddsNewKeys(t *testing.T) {
	results := []drift.Result{
		makeResult("svc-a", "image", true),
		makeResult("svc-b", "image", true),
	}
	m := stale.UpdateFirstSeen(nil, results, epoch)
	if len(m) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(m))
	}
}

func TestUpdateFirstSeen_DoesNotOverwriteExisting(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", true)}
	existing := map[string]time.Time{"svc-a:image": epoch.Add(-time.Hour)}
	m := stale.UpdateFirstSeen(existing, results, epoch)
	if !m["svc-a:image"].Equal(epoch.Add(-time.Hour)) {
		t.Error("existing key should not be overwritten")
	}
}

func TestUpdateFirstSeen_SkipsNonDrifted(t *testing.T) {
	results := []drift.Result{makeResult("svc-a", "image", false)}
	m := stale.UpdateFirstSeen(nil, results, epoch)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

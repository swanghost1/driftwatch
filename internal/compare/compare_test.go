package compare_test

import (
	"testing"

	"github.com/example/driftwatch/internal/compare"
	"github.com/example/driftwatch/internal/drift"
)

func makeDrift(service, field, want, got string) drift.Result {
	return drift.Result{
		Service: service,
		Field:   field,
		Expected: want,
		Actual:  got,
		Drifted: true,
	}
}

func makeOK(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func TestDiff_AllIntroduced(t *testing.T) {
	current := []drift.Result{makeDrift("svc-a", "image", "v1", "v2")}
	changes := compare.Diff(nil, current)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != compare.ChangeIntroduced {
		t.Errorf("expected Introduced, got %s", changes[0].Kind)
	}
}

func TestDiff_AllResolved(t *testing.T) {
	previous := []drift.Result{makeDrift("svc-a", "image", "v1", "v2")}
	changes := compare.Diff(previous, nil)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != compare.ChangeResolved {
		t.Errorf("expected Resolved, got %s", changes[0].Kind)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	r := makeDrift("svc-a", "replicas", "3", "1")
	changes := compare.Diff([]drift.Result{r}, []drift.Result{r})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != compare.ChangeUnchanged {
		t.Errorf("expected Unchanged, got %s", changes[0].Kind)
	}
}

func TestDiff_OKResultsIgnored(t *testing.T) {
	prev := []drift.Result{makeOK("svc-a")}
	curr := []drift.Result{makeOK("svc-b")}
	changes := compare.Diff(prev, curr)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes for non-drifted results, got %d", len(changes))
	}
}

func TestOnlyIntroduced(t *testing.T) {
	changes := []compare.Change{
		{Kind: compare.ChangeIntroduced},
		{Kind: compare.ChangeResolved},
		{Kind: compare.ChangeUnchanged},
	}
	got := compare.OnlyIntroduced(changes)
	if len(got) != 1 {
		t.Errorf("expected 1 introduced, got %d", len(got))
	}
}

func TestOnlyResolved(t *testing.T) {
	changes := []compare.Change{
		{Kind: compare.ChangeIntroduced},
		{Kind: compare.ChangeResolved},
		{Kind: compare.ChangeResolved},
	}
	got := compare.OnlyResolved(changes)
	if len(got) != 2 {
		t.Errorf("expected 2 resolved, got %d", len(got))
	}
}

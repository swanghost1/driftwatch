package cooldown_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/cooldown"
	"github.com/example/driftwatch/internal/drift"
)

func makeDrifted(service string) drift.Result {
	return drift.Result{Service: service, Drifted: true}
}

func makeClean(service string) drift.Result {
	return drift.Result{Service: service, Drifted: false}
}

func storeAt(t *testing.T) *cooldown.Store {
	t.Helper()
	s, err := cooldown.NewStore(filepath.Join(t.TempDir(), "cd.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestApply_CleanResults_AlwaysPass(t *testing.T) {
	s := storeAt(t)
	in := []drift.Result{makeClean("svc-a"), makeClean("svc-b")}
	out, err := cooldown.Apply(s, in, cooldown.Options{Period: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("want 2 results, got %d", len(out))
	}
}

func TestApply_FirstDrift_Passes(t *testing.T) {
	s := storeAt(t)
	in := []drift.Result{makeDrifted("svc-a")}
	out, err := cooldown.Apply(s, in, cooldown.Options{Period: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("want 1, got %d", len(out))
	}
}

func TestApply_SecondDrift_WithinPeriod_Suppressed(t *testing.T) {
	s := storeAt(t)
	in := []drift.Result{makeDrifted("svc-a")}
	_, _ = cooldown.Apply(s, in, cooldown.Options{Period: time.Hour})

	out, err := cooldown.Apply(s, in, cooldown.Options{Period: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatalf("want 0 (suppressed), got %d", len(out))
	}
}

func TestApply_SecondDrift_PeriodZero_Passes(t *testing.T) {
	s := storeAt(t)
	in := []drift.Result{makeDrifted("svc-b")}
	_, _ = cooldown.Apply(s, in, cooldown.Options{Period: 0})

	out, err := cooldown.Apply(s, in, cooldown.Options{Period: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("want 1, got %d", len(out))
	}
}

func TestApply_MixedResults_OnlyCoolingDownSuppressed(t *testing.T) {
	s := storeAt(t)
	_ = s.Record("svc-hot") // already recorded

	in := []drift.Result{makeDrifted("svc-hot"), makeDrifted("svc-new")}
	out, err := cooldown.Apply(s, in, cooldown.Options{Period: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("want 1 (svc-new only), got %d", len(out))
	}
	if out[0].Service != "svc-new" {
		t.Fatalf("expected svc-new, got %s", out[0].Service)
	}
}

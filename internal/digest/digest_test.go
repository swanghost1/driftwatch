package digest_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/digest"
)

func base() digest.ServiceSnapshot {
	return digest.ServiceSnapshot{
		Name:     "api",
		Image:    "nginx:1.25",
		Replicas: 3,
		Env:      map[string]string{"PORT": "8080", "DEBUG": "false"},
		Tags:     []string{"team:platform", "env:prod"},
	}
}

func TestCompute_ReturnsSHA256Hex(t *testing.T) {
	h, err := digest.Compute(base())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h) != 64 {
		t.Fatalf("expected 64-char hex string, got %d chars: %s", len(h), h)
	}
}

func TestCompute_Deterministic(t *testing.T) {
	h1, _ := digest.Compute(base())
	h2, _ := digest.Compute(base())
	if h1 != h2 {
		t.Fatalf("digests differ across calls: %s vs %s", h1, h2)
	}
}

func TestCompute_TagOrderIndependent(t *testing.T) {
	a := base()
	b := base()
	b.Tags = []string{"env:prod", "team:platform"} // reversed

	ha, _ := digest.Compute(a)
	hb, _ := digest.Compute(b)
	if ha != hb {
		t.Fatalf("tag order should not affect digest: %s vs %s", ha, hb)
	}
}

func TestCompute_ChangedImage_DifferentDigest(t *testing.T) {
	a := base()
	b := base()
	b.Image = "nginx:1.26"

	ha, _ := digest.Compute(a)
	hb, _ := digest.Compute(b)
	if ha == hb {
		t.Fatal("expected different digests for different images")
	}
}

func TestEqual_SameSnapshot_ReturnsTrue(t *testing.T) {
	ok, err := digest.Equal(base(), base())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Equal to return true for identical snapshots")
	}
}

func TestEqual_DifferentReplicas_ReturnsFalse(t *testing.T) {
	a := base()
	b := base()
	b.Replicas = 5

	ok, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Equal to return false for different replica counts")
	}
}

func TestEqual_DifferentEnv_ReturnsFalse(t *testing.T) {
	a := base()
	b := base()
	b.Env["PORT"] = "9090"

	ok, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Equal to return false when env differs")
	}
}

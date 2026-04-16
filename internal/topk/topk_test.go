package topk_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/topk"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "alpha", Drifted: true},
		{Service: "alpha", Drifted: true},
		{Service: "alpha", Drifted: true},
		{Service: "beta", Drifted: true},
		{Service: "beta", Drifted: true},
		{Service: "gamma", Drifted: true},
		{Service: "delta", Drifted: false},
		{Service: "delta", Drifted: false},
	}
}

func TestApply_RankedByDriftCount(t *testing.T) {
	entries := topk.Apply(makeResults(), topk.Options{K: 3})
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Service != "alpha" || entries[0].DriftCount != 3 {
		t.Errorf("expected alpha/3, got %s/%d", entries[0].Service, entries[0].DriftCount)
	}
	if entries[1].Service != "beta" || entries[1].DriftCount != 2 {
		t.Errorf("expected beta/2, got %s/%d", entries[1].Service, entries[1].DriftCount)
	}
}

func TestApply_ZeroK_ReturnsAll(t *testing.T) {
	entries := topk.Apply(makeResults(), topk.Options{K: 0})
	if len(entries) != 3 {
		t.Fatalf("expected 3 drifted services, got %d", len(entries))
	}
}

func TestApply_NoDrift_ReturnsEmpty(t *testing.T) {
	results := []drift.Result{
		{Service: "alpha", Drifted: false},
	}
	entries := topk.Apply(results, topk.DefaultOptions())
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}

func TestApply_TieBreak_Alphabetical(t *testing.T) {
	results := []drift.Result{
		{Service: "zebra", Drifted: true},
		{Service: "apple", Drifted: true},
	}
	entries := topk.Apply(results, topk.Options{K: 2})
	if entries[0].Service != "apple" {
		t.Errorf("expected apple first on tie, got %s", entries[0].Service)
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	entries := topk.Apply(makeResults(), topk.DefaultOptions())
	topk.Write(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "SERVICE") || !strings.Contains(out, "DRIFT COUNT") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestWrite_EmptyEntries_ShowsMessage(t *testing.T) {
	var buf bytes.Buffer
	topk.Write(&buf, nil)
	if !strings.Contains(buf.String(), "no drifted") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

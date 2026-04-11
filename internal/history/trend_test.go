package history_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/history"
)

func makeEntry(drifted, clean int) history.Entry {
	var results []drift.Result
	for i := 0; i < drifted; i++ {
		results = append(results, drift.Result{ServiceName: "svc", Status: drift.StatusDrifted})
	}
	for i := 0; i < clean; i++ {
		results = append(results, drift.Result{ServiceName: "svc", Status: drift.StatusOK})
	}
	return history.Entry{RunAt: time.Now().UTC(), Results: results}
}

func TestSummarise_Counts(t *testing.T) {
	entries := []history.Entry{makeEntry(3, 7)}
	summaries := history.Summarise(entries)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	s := summaries[0]
	if s.Total != 10 {
		t.Errorf("Total: got %d, want 10", s.Total)
	}
	if s.Drifted != 3 {
		t.Errorf("Drifted: got %d, want 3", s.Drifted)
	}
	if s.Clean != 7 {
		t.Errorf("Clean: got %d, want 7", s.Clean)
	}
}

func TestSummarise_DriftRate(t *testing.T) {
	entries := []history.Entry{makeEntry(1, 1)}
	summaries := history.Summarise(entries)
	if summaries[0].DriftRate != 50.0 {
		t.Errorf("DriftRate: got %.1f, want 50.0", summaries[0].DriftRate)
	}
}

func TestSummarise_EmptyEntry(t *testing.T) {
	entries := []history.Entry{{RunAt: time.Now()}}
	summaries := history.Summarise(entries)
	if summaries[0].DriftRate != 0 {
		t.Errorf("DriftRate on empty: got %.1f, want 0", summaries[0].DriftRate)
	}
}

func TestWriteTrend_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	summaries := history.Summarise([]history.Entry{makeEntry(2, 8)})
	if err := history.WriteTrend(&buf, summaries); err != nil {
		t.Fatalf("WriteTrend: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"RUN AT", "TOTAL", "DRIFTED", "CLEAN", "DRIFT RATE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
}

func TestWriteTrend_NilSummaries_WritesHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := history.WriteTrend(&buf, nil); err != nil {
		t.Fatalf("WriteTrend: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}
}

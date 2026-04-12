package scoring_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/scoring"
)

func makeEntry(label string, value float64, grade string) scoring.TrendEntry {
	return scoring.TrendEntry{
		Label: label,
		Score: scoring.Score{Value: value, Grade: grade},
	}
}

func TestDirection_Improving(t *testing.T) {
	tr := scoring.Trend{
		makeEntry("run-1", 60, "D"),
		makeEntry("run-2", 80, "B"),
		makeEntry("run-3", 95, "A"),
	}
	if tr.Direction() != "improving" {
		t.Errorf("expected improving, got %s", tr.Direction())
	}
}

func TestDirection_Degrading(t *testing.T) {
	tr := scoring.Trend{
		makeEntry("run-1", 95, "A"),
		makeEntry("run-2", 50, "D"),
	}
	if tr.Direction() != "degrading" {
		t.Errorf("expected degrading, got %s", tr.Direction())
	}
}

func TestDirection_Stable(t *testing.T) {
	tr := scoring.Trend{
		makeEntry("run-1", 80, "B"),
		makeEntry("run-2", 80, "B"),
	}
	if tr.Direction() != "stable" {
		t.Errorf("expected stable, got %s", tr.Direction())
	}
}

func TestDirection_SingleEntry(t *testing.T) {
	tr := scoring.Trend{makeEntry("run-1", 70, "C")}
	if tr.Direction() != "stable" {
		t.Errorf("single entry should be stable, got %s", tr.Direction())
	}
}

func TestWriteTrend_ContainsHeaders(t *testing.T) {
	tr := scoring.Trend{
		makeEntry("2024-01-01T00:00:00Z", 90, "A"),
		makeEntry("2024-01-02T00:00:00Z", 75, "C"),
	}
	var buf bytes.Buffer
	scoring.WriteTrend(&buf, tr)
	out := buf.String()
	for _, want := range []string{"Timestamp", "Score", "Grade", "degrading"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteTrend_Empty(t *testing.T) {
	var buf bytes.Buffer
	scoring.WriteTrend(&buf, scoring.Trend{})
	if !strings.Contains(buf.String(), "No scoring history") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

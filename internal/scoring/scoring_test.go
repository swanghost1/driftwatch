package scoring_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/scoring"
)

func makeResults(total, drifted int) []drift.Result {
	results := make([]drift.Result, total)
	for i := 0; i < drifted; i++ {
		results[i] = drift.Result{ServiceName: "svc", Drifted: true}
	}
	for i := drifted; i < total; i++ {
		results[i] = drift.Result{ServiceName: "svc", Drifted: false}
	}
	return results
}

func TestCompute_NoDrift(t *testing.T) {
	s := scoring.Compute(makeResults(4, 0))
	if s.Value != 100 {
		t.Errorf("expected 100, got %.2f", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
}

func TestCompute_AllDrifted(t *testing.T) {
	s := scoring.Compute(makeResults(4, 4))
	if s.Value != 0 {
		t.Errorf("expected 0, got %.2f", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected grade F, got %s", s.Grade)
	}
}

func TestCompute_PartialDrift(t *testing.T) {
	s := scoring.Compute(makeResults(10, 2))
	if s.Value != 80.00 {
		t.Errorf("expected 80.00, got %.2f", s.Value)
	}
	if s.Grade != "B" {
		t.Errorf("expected grade B, got %s", s.Grade)
	}
	if s.Drifted != 2 {
		t.Errorf("expected 2 drifted, got %d", s.Drifted)
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	s := scoring.Compute(nil)
	if s.Value != 100 {
		t.Errorf("expected 100 for empty results, got %.2f", s.Value)
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestGrade_Boundaries(t *testing.T) {
	cases := []struct {
		drifted, total int
		want           string
	}{
		{0, 20, "A"},
		{4, 20, "B"},  // 80%
		{7, 20, "C"},  // 65%
		{10, 20, "D"}, // 50%
		{11, 20, "F"}, // 45%
	}
	for _, tc := range cases {
		s := scoring.Compute(makeResults(tc.total, tc.drifted))
		if s.Grade != tc.want {
			t.Errorf("drifted=%d total=%d: expected grade %s, got %s", tc.drifted, tc.total, tc.want, s.Grade)
		}
	}
}

func TestWrite_ContainsScore(t *testing.T) {
	s := scoring.Compute(makeResults(5, 1))
	var buf bytes.Buffer
	scoring.Write(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "80.00") {
		t.Errorf("expected score 80.00 in output, got: %s", out)
	}
	if !strings.Contains(out, "Grade: B") {
		t.Errorf("expected Grade B in output, got: %s", out)
	}
}

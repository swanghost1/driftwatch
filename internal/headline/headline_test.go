package headline_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/headline"
)

func makeResults(drifted, clean int) []drift.Result {
	var out []drift.Result
	for i := 0; i < drifted; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < clean; i++ {
		out = append(out, drift.Result{Service: "svc", Drifted: false})
	}
	return out
}

func TestBuild_Counts(t *testing.T) {
	s := headline.Build(makeResults(3, 2))
	if s.Total != 5 || s.Drifted != 3 || s.Clean != 2 {
		t.Fatalf("unexpected counts: %+v", s)
	}
}

func TestBuild_NoDrift(t *testing.T) {
	s := headline.Build(makeResults(0, 4))
	if s.Drifted != 0 || s.Clean != 4 {
		t.Fatalf("expected all clean, got %+v", s)
	}
}

func TestBuild_Empty(t *testing.T) {
	s := headline.Build(nil)
	if s.Total != 0 {
		t.Fatalf("expected zero total")
	}
}

func TestWriteText_ContainsDriftDetected(t *testing.T) {
	s := headline.Summary{Total: 5, Drifted: 2, Clean: 3, At: time.Now()}
	var buf bytes.Buffer
	if err := headline.WriteText(&buf, s); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "DRIFT DETECTED") {
		t.Errorf("expected DRIFT DETECTED in output: %s", buf.String())
	}
}

func TestWriteText_ContainsOK(t *testing.T) {
	s := headline.Summary{Total: 3, Drifted: 0, Clean: 3, At: time.Now()}
	var buf bytes.Buffer
	_ = headline.WriteText(&buf, s)
	if !strings.Contains(buf.String(), "[OK]") {
		t.Errorf("expected [OK] in output: %s", buf.String())
	}
}

func TestWriteText_ContainsCounts(t *testing.T) {
	s := headline.Summary{Total: 4, Drifted: 1, Clean: 3, At: time.Now()}
	var buf bytes.Buffer
	_ = headline.WriteText(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "1/4") {
		t.Errorf("expected 1/4 in output: %s", out)
	}
}

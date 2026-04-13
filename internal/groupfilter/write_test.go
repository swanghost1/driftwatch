package groupfilter_test

import (
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/groupfilter"
)

func TestWriteGroupSummary_ContainsHeader(t *testing.T) {
	results := []interface{}{
		makeResult("frontend/web", false),
		makeResult("backend/api", true),
	}
	_ = results // use helper directly

	var buf strings.Builder
	groupfilter.WriteGroupSummary(&buf, []interface{}{
		makeResult("frontend/web", false),
		makeResult("backend/api", true),
	}[0:0])
	// empty input should still print header
	if !strings.Contains(buf.String(), "GROUP SUMMARY") {
		t.Error("expected GROUP SUMMARY header")
	}
}

func TestWriteGroupSummary_ShowsDriftStatus(t *testing.T) {
	var buf strings.Builder
	groupfilter.WriteGroupSummary(&buf, []interface{}{
		makeResult("backend/api", true),
		makeResult("backend/worker", false),
		makeResult("frontend/web", false),
	})
	out := buf.String()
	if !strings.Contains(out, "DRIFT") {
		t.Error("expected DRIFT in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK in output")
	}
}

func TestWriteGroupSummary_SortedAlphabetically(t *testing.T) {
	var buf strings.Builder
	groupfilter.WriteGroupSummary(&buf, []interface{}{
		makeResult("zebra/svc", false),
		makeResult("alpha/svc", true),
		makeResult("middle/svc", false),
	})
	out := buf.String()
	alphaIdx := strings.Index(out, "alpha")
	middleIdx := strings.Index(out, "middle")
	zebraIdx := strings.Index(out, "zebra")
	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Error("expected alphabetical order: alpha, middle, zebra")
	}
}

func TestWriteGroupSummary_CountsCorrect(t *testing.T) {
	var buf strings.Builder
	groupfilter.WriteGroupSummary(&buf, []interface{}{
		makeResult("team/a", true),
		makeResult("team/b", true),
		makeResult("team/c", false),
	})
	out := buf.String()
	if !strings.Contains(out, "2/3") {
		t.Errorf("expected 2/3 drifted count in output, got:\n%s", out)
	}
}

package diff_test

import (
	"strings"
	"testing"

	"github.com/your-org/driftwatch/internal/diff"
)

func TestCompare_NoDrift(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25", "replicas": "3"}
	live := map[string]string{"image": "nginx:1.25", "replicas": "3"}

	res := diff.Compare("web", declared, live)

	if res.HasDrift() {
		t.Fatalf("expected no drift, got %d field(s)", len(res.Fields))
	}
}

func TestCompare_ValueMismatch(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25"}
	live := map[string]string{"image": "nginx:1.24"}

	res := diff.Compare("web", declared, live)

	if !res.HasDrift() {
		t.Fatal("expected drift but got none")
	}
	if len(res.Fields) != 1 {
		t.Fatalf("expected 1 field diff, got %d", len(res.Fields))
	}
	f := res.Fields[0]
	if f.Name != "image" || f.Declared != "nginx:1.25" || f.Live != "nginx:1.24" {
		t.Errorf("unexpected field diff: %+v", f)
	}
}

func TestCompare_MissingLiveField(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25", "env_DEBUG": "true"}
	live := map[string]string{"image": "nginx:1.25"}

	res := diff.Compare("api", declared, live)

	if !res.HasDrift() {
		t.Fatal("expected drift for missing live field")
	}
	if res.Fields[0].Live != "<missing>" {
		t.Errorf("expected live=<missing>, got %q", res.Fields[0].Live)
	}
}

func TestCompare_ExtraLiveField(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25"}
	live := map[string]string{"image": "nginx:1.25", "env_DEBUG": "true"}

	res := diff.Compare("api", declared, live)

	if !res.HasDrift() {
		t.Fatal("expected drift for extra live field")
	}
	if res.Fields[0].Declared != "<missing>" {
		t.Errorf("expected declared=<missing>, got %q", res.Fields[0].Declared)
	}
}

func TestResult_Summary_NoDrift(t *testing.T) {
	res := diff.Result{Service: "svc"}
	if !strings.Contains(res.Summary(), "no drift") {
		t.Errorf("expected 'no drift' in summary, got %q", res.Summary())
	}
}

func TestResult_Summary_WithDrift(t *testing.T) {
	res := diff.Result{
		Service: "svc",
		Fields: []diff.Field{{Name: "image", Declared: "nginx:1.25", Live: "nginx:1.24"}},
	}
	s := res.Summary()
	if !strings.Contains(s, "svc") || !strings.Contains(s, "image") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestField_String(t *testing.T) {
	f := diff.Field{Name: "replicas", Declared: "3", Live: "2"}
	out := f.String()
	if !strings.Contains(out, "replicas") || !strings.Contains(out, "3") || !strings.Contains(out, "2") {
		t.Errorf("unexpected field string: %q", out)
	}
}

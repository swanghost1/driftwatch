package annotate_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/annotate"
	"github.com/example/driftwatch/internal/drift"
)

func TestWriteText_ContainsHeader(t *testing.T) {
	r := makeResult("payments")
	r.Annotations = map[string]string{"env": "prod"}
	var buf bytes.Buffer
	annotate.WriteText(&buf, []drift.Result{r})
	if !strings.Contains(buf.String(), "ANNOTATIONS") {
		t.Errorf("expected ANNOTATIONS header, got: %s", buf.String())
	}
}

func TestWriteText_SortedByKey(t *testing.T) {
	r := makeResult("svc")
	r.Annotations = map[string]string{"zzz": "last", "aaa": "first"}
	var buf bytes.Buffer
	annotate.WriteText(&buf, []drift.Result{r})
	aIdx := strings.Index(buf.String(), "aaa")
	zIdx := strings.Index(buf.String(), "zzz")
	if aIdx > zIdx {
		t.Errorf("expected aaa before zzz in output")
	}
}

func TestWriteJSON_ContainsServiceField(t *testing.T) {
	r := makeResult("inventory")
	r.Annotations = map[string]string{"k": "v"}
	var buf bytes.Buffer
	if err := annotate.WriteJSON(&buf, []drift.Result{r}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0]["service"] != "inventory" {
		t.Errorf("expected service=inventory, got %v", out[0]["service"])
	}
}

func TestWriteJSON_EmptyAnnotations_StillValid(t *testing.T) {
	var buf bytes.Buffer
	if err := annotate.WriteJSON(&buf, []drift.Result{makeResult("svc")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !json.Valid(buf.Bytes()) {
		t.Errorf("output is not valid JSON: %s", buf.String())
	}
}

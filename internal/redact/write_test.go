package redact_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/redact"
)

func TestWriteText_NoRedactedFields_ShowsNone(t *testing.T) {
	results := []redact.Result{
		{Service: "svc", Field: "image", Expected: "nginx:1", Actual: "nginx:2", Drifted: true},
	}
	var buf bytes.Buffer
	redact.WriteText(&buf, results)
	if !strings.Contains(buf.String(), "no sensitive fields redacted") {
		t.Errorf("expected no-redaction message, got: %s", buf.String())
	}
}

func TestWriteText_RedactedFields_ListsThem(t *testing.T) {
	results := redact.Apply(makeResults(), redact.Options{})
	var buf bytes.Buffer
	redact.WriteText(&buf, results)
	if !strings.Contains(buf.String(), "[redacted]") {
		t.Errorf("expected redacted marker in output")
	}
	if !strings.Contains(buf.String(), "total redacted fields") {
		t.Errorf("expected total count in output")
	}
}

func TestWriteText_CountsCorrectly(t *testing.T) {
	results := redact.Apply(makeResults(), redact.Options{})
	var buf bytes.Buffer
	redact.WriteText(&buf, results)
	// PASSWORD and API_KEY should both be redacted
	if !strings.Contains(buf.String(), "total redacted fields: 2") {
		t.Errorf("expected 2 redacted fields, output: %s", buf.String())
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	results := redact.Apply(makeResults(), redact.Options{})
	var buf bytes.Buffer
	if err := redact.WriteJSON(&buf, results); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var decoded []redact.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(decoded))
	}
}

func TestWriteJSON_SensitiveValuesRedacted(t *testing.T) {
	results := redact.Apply(makeResults(), redact.Options{})
	var buf bytes.Buffer
	_ = redact.WriteJSON(&buf, results)
	if strings.Contains(buf.String(), "hunter2") || strings.Contains(buf.String(), "letmein") {
		t.Error("raw sensitive values should not appear in JSON output")
	}
}

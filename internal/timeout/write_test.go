package timeout_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/timeout"
)

func makeRecord(timedOut bool) timeout.Record {
	return timeout.Record{
		Deadline:    30 * time.Second,
		GracePeriod: 5 * time.Second,
		Elapsed:     12500 * time.Millisecond,
		TimedOut:    timedOut,
	}
}

func TestWriteText_OKStatus(t *testing.T) {
	var buf bytes.Buffer
	if err := timeout.WriteText(&buf, makeRecord(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ok") {
		t.Errorf("expected 'ok' in output, got: %s", buf.String())
	}
}

func TestWriteText_TimedOutStatus(t *testing.T) {
	var buf bytes.Buffer
	if err := timeout.WriteText(&buf, makeRecord(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "TIMED OUT") {
		t.Errorf("expected 'TIMED OUT' in output, got: %s", buf.String())
	}
}

func TestWriteText_ContainsElapsed(t *testing.T) {
	var buf bytes.Buffer
	if err := timeout.WriteText(&buf, makeRecord(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "elapsed") {
		t.Errorf("expected 'elapsed' in output, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := timeout.WriteJSON(&buf, makeRecord(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteJSON_TimedOutField(t *testing.T) {
	var buf bytes.Buffer
	if err := timeout.WriteJSON(&buf, makeRecord(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "timed_out") {
		t.Errorf("expected 'timed_out' key in JSON, got: %s", buf.String())
	}
}

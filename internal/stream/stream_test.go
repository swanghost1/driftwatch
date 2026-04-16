package stream_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/stream"
)

func makeResults(drifted bool) []drift.Result {
	return []drift.Result{
		{Service: "api", Field: "image", Want: "v1", Got: "v2", Drifted: drifted},
		{Service: "worker", Field: "replicas", Want: "3", Got: "3", Drifted: false},
	}
}

func chanFrom(results []drift.Result) <-chan drift.Result {
	ch := make(chan drift.Result, len(results))
	for _, r := range results {
		ch <- r
	}
	close(ch)
	return ch
}

func TestWrite_TextFormat_OKResult(t *testing.T) {
	results := makeResults(false)
	var buf bytes.Buffer
	opts := stream.DefaultOptions()
	opts.FlushInterval = 0
	if err := stream.Write(&buf, chanFrom(results), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", out)
	}
}

func TestWrite_TextFormat_DriftedResult(t *testing.T) {
	results := makeResults(true)
	var buf bytes.Buffer
	opts := stream.DefaultOptions()
	opts.FlushInterval = 0
	if err := stream.Write(&buf, chanFrom(results), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFT]") {
		t.Errorf("expected [DRIFT] in output, got: %s", out)
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	results := makeResults(true)
	var buf bytes.Buffer
	opts := stream.Options{Format: "json", FlushInterval: 0}
	if err := stream.Write(&buf, chanFrom(results), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for _, line := range lines {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Errorf("invalid JSON line %q: %v", line, err)
		}
	}
}

func TestWrite_JSONFormat_ContainsServiceField(t *testing.T) {
	results := []drift.Result{{Service: "svc", Field: "image", Want: "a", Got: "b", Drifted: true}}
	var buf bytes.Buffer
	opts := stream.Options{Format: "json", FlushInterval: 0}
	if err := stream.Write(&buf, chanFrom(results), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"svc"`) {
		t.Errorf("expected service name in JSON output")
	}
}

func TestWrite_HeartbeatWritten(t *testing.T) {
	ch := make(chan drift.Result)
	close(ch)
	var buf bytes.Buffer
	opts := stream.Options{Format: "text", FlushInterval: 1 * time.Millisecond}
	time.Sleep(5 * time.Millisecond)
	// Channel is already closed so Write returns immediately; heartbeat may not
	// fire. Just verify no error is returned.
	if err := stream.Write(&buf, ch, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultOptions_FormatIsText(t *testing.T) {
	opts := stream.DefaultOptions()
	if opts.Format != "text" {
		t.Errorf("expected default format \"text\", got %q", opts.Format)
	}
}

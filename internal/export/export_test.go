package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/export"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Drifts: []drift.Drift{
				{Field: "image", Declared: "nginx:1.25", Live: "nginx:1.24"},
			},
		},
		{
			Service: "worker",
			Drifts:  []drift.Drift{},
		},
	}
}

func TestWrite_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, makeResults(), export.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "service,status,field,declared,live") {
		t.Errorf("expected CSV header, got:\n%s", buf.String())
	}
}

func TestWrite_CSV_DriftRow(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, makeResults(), export.FormatCSV)
	if !strings.Contains(buf.String(), "api,drift,image,nginx:1.25,nginx:1.24") {
		t.Errorf("expected drift row, got:\n%s", buf.String())
	}
}

func TestWrite_CSV_OKRow(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, makeResults(), export.FormatCSV)
	if !strings.Contains(buf.String(), "worker,ok") {
		t.Errorf("expected ok row, got:\n%s", buf.String())
	}
}

func TestWrite_Markdown_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, makeResults(), export.FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "| Service |") {
		t.Errorf("expected markdown header, got:\n%s", buf.String())
	}
}

func TestWrite_Markdown_DriftRow(t *testing.T) {
	var buf bytes.Buffer
	_ = export.Write(&buf, makeResults(), export.FormatMarkdown)
	if !strings.Contains(buf.String(), "drift") {
		t.Errorf("expected drift in markdown output, got:\n%s", buf.String())
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, makeResults(), export.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

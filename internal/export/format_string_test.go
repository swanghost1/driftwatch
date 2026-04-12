package export_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/export"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  export.Format
	}{
		{"csv", export.FormatCSV},
		{"markdown", export.FormatMarkdown},
	}
	for _, tc := range cases {
		got, err := export.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := export.ParseFormat("toml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormat_String(t *testing.T) {
	if export.FormatCSV.String() != "csv" {
		t.Errorf("FormatCSV.String() = %q, want \"csv\"", export.FormatCSV.String())
	}
	if export.FormatMarkdown.String() != "markdown" {
		t.Errorf("FormatMarkdown.String() = %q, want \"markdown\"", export.FormatMarkdown.String())
	}
}

func TestKnownFormats_Length(t *testing.T) {
	if len(export.KnownFormats()) != 2 {
		t.Errorf("expected 2 known formats, got %d", len(export.KnownFormats()))
	}
}

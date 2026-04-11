package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/output"
)

func TestOpen_EmptyPath_ReturnsStdout(t *testing.T) {
	wc, err := output.Open(output.Destination{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer wc.Close()

	// Writing to stdout should not error.
	_, err = wc.Write([]byte(""))
	if err != nil {
		t.Fatalf("write to stdout wrapper failed: %v", err)
	}
}

func TestOpen_FilePath_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out", "results.txt")

	wc, err := output.Open(output.Destination{Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, _ = wc.Write([]byte("hello"))
	wc.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(data))
	}
}

func TestOpen_Truncates_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	_ = os.WriteFile(path, []byte(strings.Repeat("x", 100)), 0o644)

	wc, err := output.Open(output.Destination{Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _ = wc.Write([]byte("new"))
	wc.Close()

	data, _ := os.ReadFile(path)
	if string(data) != "new" {
		t.Errorf("expected truncated content %q, got %q", "new", string(data))
	}
}

func TestOpen_Append_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	_ = os.WriteFile(path, []byte("first"), 0o644)

	wc, err := output.Open(output.Destination{Path: path, Append: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _ = wc.Write([]byte("-second"))
	wc.Close()

	data, _ := os.ReadFile(path)
	if string(data) != "first-second" {
		t.Errorf("expected appended content, got %q", string(data))
	}
}

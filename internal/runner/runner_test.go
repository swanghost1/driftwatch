package runner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/runner"
)

const validConfig = `version: "1"
services:
  - name: api
    image: nginx:1.25
    replicas: 2
    env:
      PORT: "8080"
`

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestRun_NoConfigFile(t *testing.T) {
	code, err := runner.Run(runner.Options{
		ConfigPath: "/nonexistent/path.yaml",
		Format:     "text",
	})
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_ValidConfig_NoDrift(t *testing.T) {
	path := writeTempConfig(t, validConfig)
	code, err := runner.Run(runner.Options{
		ConfigPath:  path,
		Format:      "text",
		FailOnDrift: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No drift expected — even with FailOnDrift the code should be 0.
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_OutputToFile(t *testing.T) {
	configPath := writeTempConfig(t, validConfig)
	outPath := filepath.Join(t.TempDir(), "report.txt")

	code, err := runner.Run(runner.Options{
		ConfigPath: configPath,
		Format:     "text",
		OutputPath: outPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "api") {
		t.Errorf("expected output to mention service 'api', got:\n%s", string(data))
	}
}

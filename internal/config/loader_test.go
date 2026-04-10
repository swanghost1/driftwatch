package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `
version: "1"
services:
  - name: api
    image: myrepo/api:1.2.3
    environment: production
    replicas: 3
    ports: [8080]
    env_vars:
      LOG_LEVEL: info
`
	path := writeTempConfig(t, content)
	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	svc := cfg.Services[0]
	if svc.Name != "api" {
		t.Errorf("expected name 'api', got %q", svc.Name)
	}
	if svc.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", svc.Replicas)
	}
	if svc.EnvVars["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", svc.EnvVars["LOG_LEVEL"])
	}
}

func TestLoadFromFile_MissingVersion(t *testing.T) {
	content := `
services:
  - name: api
    image: myrepo/api:latest
`
	path := writeTempConfig(t, content)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestLoadFromFile_MissingImage(t *testing.T) {
	content := `
version: "1"
services:
  - name: worker
`
	path := writeTempConfig(t, content)
	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for missing image, got nil")
	}
}

func TestLoadFromFile_FileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/driftwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

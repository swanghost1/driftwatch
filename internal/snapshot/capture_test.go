package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/snapshot"
)

func baseConfig() *config.Config {
	return &config.Config{
		Version: "1",
		Services: []config.Service{
			{
				Name:     "api",
				Image:    "api:v2.0.0",
				Replicas: 2,
				Env:      map[string]string{"PORT": "8080"},
			},
			{
				Name:     "cache",
				Image:    "redis:7",
				Replicas: 1,
			},
		},
	}
}

func TestCaptureFromConfig_SnapshotCount(t *testing.T) {
	cfg := baseConfig()
	store, err := snapshot.CaptureFromConfig(cfg, snapshot.CaptureOptions{})
	if err != nil {
		t.Fatalf("CaptureFromConfig: %v", err)
	}
	if len(store.Snapshots) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(store.Snapshots))
	}
}

func TestCaptureFromConfig_FieldsPreserved(t *testing.T) {
	cfg := baseConfig()
	store, err := snapshot.CaptureFromConfig(cfg, snapshot.CaptureOptions{})
	if err != nil {
		t.Fatalf("CaptureFromConfig: %v", err)
	}

	api := store.FindByName("api")
	if api == nil {
		t.Fatal("api snapshot not found")
	}
	if api.Image != "api:v2.0.0" {
		t.Errorf("image: got %q", api.Image)
	}
	if api.Replicas != 2 {
		t.Errorf("replicas: got %d", api.Replicas)
	}
	if api"PORT"] != "8080" {
		t.Errorf("env PORT: got %q", api.Env["PORT"])
	}
	if api.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestCaptureFromConfig_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.")

	_, err := snapshot.CaptureFromConfig(baseConfig(), snapshot.CaptureOptions{OutputPath: path})
	if err != nil {
		t.Fatalf("CaptureFromConfig:t}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected snapshot file to exist: %v", err)
	}
}

func TestCaptureFromConfig_NConfig(t *testing.T) {
	_, err := snapshot.CaptureFromConfig(nil, snapshot.CaptureOptions{})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

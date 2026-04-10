package drift

import (
	"strings"
	"testing"
)

func baseService() (DeclaredService, ServiceState) {
	declared := DeclaredService{
		Name:     "api",
		Image:    "myrepo/api:v1.2.3",
		Replicas: 3,
		Env:      map[string]string{"LOG_LEVEL": "info", "PORT": "8080"},
	}
	live := ServiceState{
		Name:     "api",
		Image:    "myrepo/api:v1.2.3",
		Replicas: 3,
		Env:      map[string]string{"LOG_LEVEL": "info", "PORT": "8080"},
	}
	return declared, live
}

func TestDetect_NoDrift(t *testing.T) {
	declared, live := baseService()
	result := Detect(declared, live)
	if result.Drifted {
		t.Errorf("expected no drift, got messages: %v", result.Messages)
	}
}

func TestDetect_ImageMismatch(t *testing.T) {
	declared, live := baseService()
	live.Image = "myrepo/api:v1.2.4"
	result := Detect(declared, live)
	if !result.Drifted {
		t.Fatal("expected drift due to image mismatch")
	}
	if !containsSubstring(result.Messages, "image mismatch") {
		t.Errorf("expected image mismatch message, got: %v", result.Messages)
	}
}

func TestDetect_ReplicasMismatch(t *testing.T) {
	declared, live := baseService()
	live.Replicas = 1
	result := Detect(declared, live)
	if !result.Drifted {
		t.Fatal("expected drift due to replicas mismatch")
	}
	if !containsSubstring(result.Messages, "replicas mismatch") {
		t.Errorf("expected replicas mismatch message, got: %v", result.Messages)
	}
}

func TestDetect_EnvVarMissing(t *testing.T) {
	declared, live := baseService()
	delete(live.Env, "PORT")
	result := Detect(declared, live)
	if !result.Drifted {
		t.Fatal("expected drift due to missing env var")
	}
	if !containsSubstring(result.Messages, "PORT") {
		t.Errorf("expected PORT missing message, got: %v", result.Messages)
	}
}

func TestDetect_EnvVarValueMismatch(t *testing.T) {
	declared, live := baseService()
	live.Env["LOG_LEVEL"] = "debug"
	result := Detect(declared, live)
	if !result.Drifted {
		t.Fatal("expected drift due to env var value mismatch")
	}
	if !containsSubstring(result.Messages, "LOG_LEVEL") {
		t.Errorf("expected LOG_LEVEL mismatch message, got: %v", result.Messages)
	}
}

func TestDriftResult_Summary_NoDrift(t *testing.T) {
	r := DriftResult{Service: "svc", Drifted: false}
	if !strings.Contains(r.Summary(), "[OK]") {
		t.Errorf("expected [OK] in summary, got: %s", r.Summary())
	}
}

func TestDriftResult_Summary_Drifted(t *testing.T) {
	r := DriftResult{Service: "svc", Drifted: true, Messages: []string{"image mismatch"}}
	if !strings.Contains(r.Summary(), "[DRIFT]") {
		t.Errorf("expected [DRIFT] in summary, got: %s", r.Summary())
	}
}

func containsSubstring(messages []string, sub string) bool {
	for _, m := range messages {
		if strings.Contains(m, sub) {
			return true
		}
	}
	return false
}

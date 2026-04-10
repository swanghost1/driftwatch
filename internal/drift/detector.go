package drift

import (
	"fmt"
	"strings"
)

// ServiceState represents the live state of a deployed service.
type ServiceState struct {
	Name    string
	Image   string
	Replicas int
	Env     map[string]string
}

// DriftResult holds the comparison outcome for a single service.
type DriftResult struct {
	Service  string
	Drifted  bool
	Messages []string
}

// DeclaredService mirrors the relevant fields from the config for comparison.
type DeclaredService struct {
	Name     string
	Image    string
	Replicas int
	Env      map[string]string
}

// Detect compares a declared service config against its live state and
// returns a DriftResult describing any discrepancies found.
func Detect(declared DeclaredService, live ServiceState) DriftResult {
	result := DriftResult{
		Service: declared.Name,
	}

	if declared.Image != live.Image {
		result.Messages = append(result.Messages,
			fmt.Sprintf("image mismatch: declared=%q live=%q", declared.Image, live.Image))
	}

	if declared.Replicas != 0 && declared.Replicas != live.Replicas {
		result.Messages = append(result.Messages,
			fmt.Sprintf("replicas mismatch: declared=%d live=%d", declared.Replicas, live.Replicas))
	}

	for key, declaredVal := range declared.Env {
		liveVal, ok := live.Env[key]
		if !ok {
			result.Messages = append(result.Messages,
				fmt.Sprintf("env var %q missing in live service", key))
		} else if declaredVal != liveVal {
			result.Messages = append(result.Messages,
				fmt.Sprintf("env var %q mismatch: declared=%q live=%q", key, declaredVal, liveVal))
		}
	}

	result.Drifted = len(result.Messages) > 0
	return result
}

// Summary returns a human-readable summary of the drift result.
func (r DriftResult) Summary() string {
	if !r.Drifted {
		return fmt.Sprintf("[OK] %s: no drift detected", r.Service)
	}
	return fmt.Sprintf("[DRIFT] %s:\n  - %s", r.Service, strings.Join(r.Messages, "\n  - "))
}

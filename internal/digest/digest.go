// Package digest computes and compares deterministic hashes of service
// snapshots, allowing driftwatch to quickly detect whether a service's
// declared state has changed since the last run without performing a full
// field-by-field comparison.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

// ServiceSnapshot is the minimal view of a service used for hashing.
type ServiceSnapshot struct {
	Name     string            `json:"name"`
	Image    string            `json:"image"`
	Replicas int               `json:"replicas"`
	Env      map[string]string `json:"env,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
}

// Compute returns a stable SHA-256 hex digest for the given snapshot.
// Map keys and slice elements are sorted before hashing so that the digest
// is independent of insertion order.
func Compute(s ServiceSnapshot) (string, error) {
	normalised := struct {
		Name     string            `json:"name"`
		Image    string            `json:"image"`
		Replicas int               `json:"replicas"`
		Env      map[string]string `json:"env,omitempty"`
		Tags     []string          `json:"tags,omitempty"`
	}{
		Name:     s.Name,
		Image:    s.Image,
		Replicas: s.Replicas,
		Env:      s.Env,
	}

	if len(s.Tags) > 0 {
		copy := append([]string(nil), s.Tags...)
		sort.Strings(copy)
		normalised.Tags = copy
	}

	b, err := json.Marshal(normalised)
	if err != nil {
		return "", fmt.Errorf("digest: marshal: %w", err)
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Equal returns true when both snapshots produce the same digest.
// The error from either Compute call is returned if hashing fails.
func Equal(a, b ServiceSnapshot) (bool, error) {
	ha, err := Compute(a)
	if err != nil {
		return false, err
	}
	hb, err := Compute(b)
	if err != nil {
		return false, err
	}
	return ha == hb, nil
}

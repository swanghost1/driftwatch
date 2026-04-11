package snapshot

import (
	"fmt"
	"time"

	"github.com/example/driftwatch/internal/config"
)

// CaptureOptions controls how a snapshot is taken.
type CaptureOptions struct {
	// OutputPath is where the snapshot JSON file will be written.
	OutputPath string
}

// CaptureFromConfig builds a Store from the declared config and saves it.
// In a real implementation this would query live infrastructure; here it
// records the declared state so it can be diffed against future runs.
func CaptureFromConfig(cfg *config.Config, opts CaptureOptions) (*Store, error) {
	if cfg == nil {
		return nil, fmt.Errorf("capture: config must not be nil")
	}

	store := &Store{
		Version:   cfg.Version,
		Snapshots: make([]ServiceSnapshot, 0, len(cfg.Services)),
	}

	for _, svc := range cfg.Services {
		snap := ServiceSnapshot{
			Name:       svc.Name,
			Image:      svc.Image,
			Replicas:   svc.Replicas,
			CapturedAt: time.Now().UTC(),
		}
		if len(svc.Env) > 0 {
			snap.Env = make(map[string]string, len(svc.Env))
			for k, v := range svc.Env {
				snap.Env[k] = v
			}
		}
		store.Snapshots = append(store.Snapshots, snap)
	}

	if opts.OutputPath != "" {
		if err := Save(opts.OutputPath, store); err != nil {
			return nil, fmt.Errorf("capture: save snapshot: %w", err)
		}
	}

	return store, nil
}

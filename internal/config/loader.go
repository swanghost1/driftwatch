package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the declared state of a service.
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Environment string            `yaml:"environment"`
	EnvVars     map[string]string `yaml:"env_vars"`
	Replicas    int               `yaml:"replicas"`
	Image       string            `yaml:"image"`
	Ports       []int             `yaml:"ports"`
}

// DriftConfig is the top-level structure of a driftwatch config file.
type DriftConfig struct {
	Version  string          `yaml:"version"`
	Services []ServiceConfig `yaml:"services"`
}

// LoadFromFile reads and parses a YAML config file at the given path.
func LoadFromFile(path string) (*DriftConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg DriftConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// ServiceByName returns the ServiceConfig with the given name, or an error if
// no service with that name exists in the config.
func (c *DriftConfig) ServiceByName(name string) (*ServiceConfig, error) {
	for i := range c.Services {
		if c.Services[i].Name == name {
			return &c.Services[i], nil
		}
	}
	return nil, fmt.Errorf("service %q not found in config", name)
}

// validate performs basic sanity checks on the loaded config.
func validate(cfg *DriftConfig) error {
	if cfg.Version == "" {
		return fmt.Errorf("config version is required")
	}
	if len(cfg.Services) == 0 {
		return fmt.Errorf("at least one service must be declared")
	}
	for i, svc := range cfg.Services {
		if svc.Name == "" {
			return fmt.Errorf("service at index %d is missing a name", i)
		}
		if svc.Image == "" {
			return fmt.Errorf("service %q is missing an image", svc.Name)
		}
	}
	return nil
}

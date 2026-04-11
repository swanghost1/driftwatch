package runner

import (
	"fmt"
	"os"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/report"
)

// Options holds the runtime options passed from the CLI.
type Options struct {
	ConfigPath string
	Format     string
	OutputPath string
	FailOnDrift bool
}

// Run loads the config, runs drift detection, and writes the report.
// It returns a non-zero exit code when drift is detected and FailOnDrift is set.
func Run(opts Options) (int, error) {
	cfg, err := config.LoadFromFile(opts.ConfigPath)
	if err != nil {
		return 1, fmt.Errorf("loading config: %w", err)
	}

	results := drift.Detect(cfg.Services)

	summary := report.NewSummary(results)

	var out *os.File
	if opts.OutputPath == "" || opts.OutputPath == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(opts.OutputPath)
		if err != nil {
			return 1, fmt.Errorf("opening output file: %w", err)
		}
		defer out.Close()
	}

	if err := report.Write(summary, opts.Format, out); err != nil {
		return 1, fmt.Errorf("writing report: %w", err)
	}

	if opts.FailOnDrift && summary.DriftCount > 0 {
		return 2, nil
	}

	return 0, nil
}

package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/cjonesy/doctor/internal/check"
	"github.com/cjonesy/doctor/internal/errors"
	"github.com/cjonesy/doctor/internal/logger"
	"gopkg.in/yaml.v3"
)

// Config is the configuration typically found in a repo's .doctor.yml file
type Config struct {
	Path        string
	Verbose     bool
	Checks      []check.Check
	Timeout     time.Duration
	Logger      *logger.Logger
	Parallelism int // Number of concurrent checks (0 = sequential, -1 = unlimited)
}

// findConfig attempts to locate the config file by walking up directories
func (c *Config) findConfig(ctx context.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Walk up directories to see if we can find a config
	for dir != "/" && err == nil {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("config search cancelled: %w", ctx.Err())
		default:
		}

		cfgPath := filepath.Join(dir, ".doctor.yml")
		if _, err := os.Stat(cfgPath); err == nil {
			c.Logger.DebugContext(ctx, "found config", "path", cfgPath)
			if c.Verbose {
				fmt.Printf("Found config at: %s\n", cfgPath)
			}
			c.Path = cfgPath
			return nil
		}

		// Drop the end of the path
		dir = path.Dir(dir)
	}

	return errors.ErrConfigNotFound
}

// parseConfig parses to config and updates the config struct
func (c *Config) parseConfig(ctx context.Context) error {
	cleanPath := filepath.Clean(c.Path)

	b, err := os.ReadFile(cleanPath)
	if err != nil {
		return errors.NewConfigError(cleanPath, fmt.Errorf("failed to read: %w", err))
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return errors.NewConfigError(cleanPath, fmt.Errorf("failed to parse YAML: %w", err))
	}

	c.Logger.DebugContext(ctx, "parsed config", "checks", len(c.Checks))
	return nil
}

// runChecks runs each check in the config (sequential or concurrent)
func (c *Config) runChecks(ctx context.Context) error {
	if c.Parallelism == 0 {
		// Sequential execution (original behavior)
		return c.runChecksSequential(ctx)
	}

	// Concurrent execution
	return c.runChecksParallel(ctx)
}

// runChecksSequential runs checks one at a time
func (c *Config) runChecksSequential(ctx context.Context) error {
	var firstErr error

	for _, chk := range c.Checks {
		select {
		case <-ctx.Done():
			return fmt.Errorf("check execution cancelled: %w", ctx.Err())
		default:
		}

		chk.Verbose = c.Verbose
		chk.Logger = c.Logger
		chk.Timeout = c.Timeout

		if err := chk.Run(ctx); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			c.Logger.ErrorContext(ctx, "check error", "error", err)
		}
	}

	return firstErr
}

// runChecksParallel runs checks concurrently using errgroup
func (c *Config) runChecksParallel(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Set concurrency limit if specified
	if c.Parallelism > 0 {
		g.SetLimit(c.Parallelism)
	}

	c.Logger.InfoContext(ctx, "running checks in parallel",
		"total", len(c.Checks),
		"limit", c.Parallelism,
	)

	// Launch goroutine for each check
	for i := range c.Checks {
		// Capture loop variable
		chk := c.Checks[i]

		g.Go(func() error {
			chk.Verbose = c.Verbose
			chk.Logger = c.Logger
			chk.Timeout = c.Timeout

			if err := chk.Run(ctx); err != nil {
				c.Logger.ErrorContext(ctx, "check error", "error", err)
				// Continue running other checks even on error
				return err
			}

			return nil
		})
	}

	// Wait for all checks to complete
	// errgroup returns first non-nil error
	if err := g.Wait(); err != nil {
		return fmt.Errorf("one or more checks failed: %w", err)
	}

	return nil
}

// Run runs a config and all checks
func (c *Config) Run(ctx context.Context) error {
	if c.Path == "" {
		if err := c.findConfig(ctx); err != nil {
			return fmt.Errorf("config discovery failed: %w", err)
		}
	}

	if err := c.parseConfig(ctx); err != nil {
		return fmt.Errorf("config parsing failed: %w", err)
	}

	if err := c.runChecks(ctx); err != nil {
		return fmt.Errorf("check execution failed: %w", err)
	}

	return nil
}

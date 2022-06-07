package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cjonesy/doctor/internal/check"
	"gopkg.in/yaml.v3"
)

// Config is the configuration typically found in a repo's .doctor.yml file
type Config struct {
	// The path in which the config lives, if set to "" we'll attempt to locate it
	Path string
	// Enables verbose output
	Verbose bool
	// A list of checks for this config
	Checks []check.Check
}

// findConfig attempts to locate the config file by walking up directories
func (c *Config) findConfig() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Walk up directories to see if we can find a config
	for dir != "/" && err == nil {
		cfgPath := filepath.Join(dir, ".doctor.yml")
		if _, err := os.Stat(cfgPath); err == nil {
			if c.Verbose {
				fmt.Printf("Found config at: %s", cfgPath)
			}
			c.Path = cfgPath
			return nil
		}

		// Drop the end of the path
		dir = path.Dir(dir)
	}

	return fmt.Errorf("config could not be found")
}

// parseConfig parses to config and updates the config struct
func (c *Config) parseConfig() error {
	cleanPath := filepath.Clean(c.Path)

	b, err := ioutil.ReadFile(cleanPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return err
	}

	return nil
}

// runChecks runs each check in the config
func (c *Config) runChecks() {
	for _, chk := range c.Checks {
		chk.Verbose = c.Verbose
		if err := chk.Run(); err != nil {
			fmt.Printf("Unexpected error:\n%s\n", err)
		}
	}
}

// Run runs a config and all checks
func (c *Config) Run() error {
	if c.Path == "" {
		if err := c.findConfig(); err != nil {
			return err
		}
	}

	if err := c.parseConfig(); err != nil {
		return err
	}

	c.runChecks()

	return nil
}

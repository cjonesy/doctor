package check

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cjonesy/doctor/internal/logger"
	"github.com/cjonesy/doctor/internal/util"
)

// Check is used to store the configuration for a check
type Check struct {
	// A description of this check
	Description string
	// The type of check to perform
	Type string
	// Depending on context this is either the command to run, or command to check
	Command string
	// Instructions the user can follow to correct an issue discovered by a check
	Fix string
	// The path in which a file exists, or is expected to exist
	Path string
	// The text that is expected
	Content string
	// Enables verbose output of checks
	Verbose bool
	// Logger for structured logging
	Logger *logger.Logger
	// Timeout for check execution
	Timeout time.Duration
}

// Run executes the configured check with context
func (c *Check) Run(ctx context.Context) error {
	// Apply timeout if set
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	c.Logger.CheckStarted(ctx, c.Type, c.Description)
	util.PrintDescription(ctx, c.Logger, c.Description)

	switch t := c.Type; t {
	case "command-in-path":
		return c.checkCommandInPath(ctx)
	case "file-exists":
		return c.checkFileExists(ctx)
	case "file-contains":
		return c.checkFileContains(ctx)
	case "output-contains":
		return c.checkOutputContains(ctx)
	default:
		util.PrintFailed()
		util.PrintFix(fmt.Sprintf("Check config file and update type '%s' to a valid type.", c.Type))
		err := fmt.Errorf("unknown check type: '%s'", c.Type)
		c.Logger.CheckError(ctx, c.Type, c.Description, err)
		return nil
	}
}

// checkCommandInPath checks that the configured command is in the user's path
func (c *Check) checkCommandInPath(ctx context.Context) error {
	path, err := exec.LookPath(c.Command)

	if err == nil {
		util.PrintPassed()
		c.Logger.CheckPassed(ctx, c.Type, c.Description)
		if c.Verbose {
			fmt.Printf("Found %s at: %s\n", c.Command, path)
		}
		return nil
	}

	util.PrintFailed()
	util.PrintFix(c.Fix)
	c.Logger.CheckFailed(ctx, c.Type, c.Description, c.Fix)
	if c.Verbose {
		fmt.Printf("Error: %v\n", err)
	}

	return nil
}

// checkFileExists checks that the configured file exists on the user's system
func (c *Check) checkFileExists(ctx context.Context) error {
	cleanPath, err := util.CleanHome(c.Path)
	if err != nil {
		util.PrintFailed()
		return fmt.Errorf("failed to clean path '%s': %w", c.Path, err)
	}

	_, err = os.Stat(cleanPath)
	if err == nil {
		util.PrintPassed()
		c.Logger.CheckPassed(ctx, c.Type, c.Description)
		if c.Verbose {
			fmt.Printf("Found %s!\n", cleanPath)
		}
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		util.PrintFailed()
		util.PrintFix(c.Fix)
		c.Logger.CheckFailed(ctx, c.Type, c.Description, c.Fix)
		if c.Verbose {
			fmt.Printf("File not found: %s\n", cleanPath)
		}
		return nil
	}

	util.PrintFailed()
	return fmt.Errorf("failed to stat file '%s': %w", cleanPath, err)
}

// checkFileContains checks that the configured file contains specific text
func (c *Check) checkFileContains(ctx context.Context) error {
	cleanPath, err := util.CleanHome(c.Path)
	if err != nil {
		util.PrintFailed()
		return fmt.Errorf("failed to clean path '%s': %w", c.Path, err)
	}

	f, err := os.ReadFile(cleanPath)
	if err != nil {
		util.PrintFailed()
		return fmt.Errorf("failed to read file '%s': %w", cleanPath, err)
	}

	found := strings.Contains(string(f), c.Content)

	if found {
		util.PrintPassed()
		c.Logger.CheckPassed(ctx, c.Type, c.Description)
		if c.Verbose {
			fmt.Printf("Found '%s' in %s!\n", c.Content, cleanPath)
		}
		return nil
	}

	util.PrintFailed()
	util.PrintFix(c.Fix)
	c.Logger.CheckFailed(ctx, c.Type, c.Description, c.Fix)
	if c.Verbose {
		fmt.Printf("Could not find '%s' in %s!\n", c.Content, cleanPath)
	}

	return nil
}

// checkOutputContains checks that the configured command outputs specific text
func (c *Check) checkOutputContains(ctx context.Context) error {
	// Use CommandContext for timeout support
	cmd := exec.CommandContext(ctx, "bash", "-c", c.Command)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		// Check if it was a timeout
		if errors.Is(err, context.DeadlineExceeded) {
			util.PrintFailed()
			util.PrintFix(fmt.Sprintf("Command timed out after %s", c.Timeout))
			c.Logger.CheckError(ctx, c.Type, c.Description, err)
			return fmt.Errorf("command timeout: %w", err)
		}

		util.PrintFailed()
		c.Logger.CheckError(ctx, c.Type, c.Description, err)
		return fmt.Errorf("command failed: %w", err)
	}

	found := strings.Contains(cmdOutput.String(), c.Content)

	if found {
		util.PrintPassed()
		c.Logger.CheckPassed(ctx, c.Type, c.Description)
		if c.Verbose {
			fmt.Printf("Found '%s' in '%s' output!\n", c.Content, c.Command)
			fmt.Print(cmdOutput.String())
		}
		return nil
	}

	util.PrintFailed()
	util.PrintFix(c.Fix)
	c.Logger.CheckFailed(ctx, c.Type, c.Description, c.Fix)
	if c.Verbose {
		fmt.Printf("Could not find '%s' in '%s' output!\n", c.Content, c.Command)
		fmt.Print(cmdOutput.String())
	}

	return nil
}

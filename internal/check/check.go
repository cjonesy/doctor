package check

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
}

// Run executes the configured check
func (c *Check) Run() error {
	util.PrintDescription(c.Description)

	switch t := c.Type; t {
	case "command-in-path":
		return c.checkCommandInPath()
	case "file-exists":
		return c.checkFileExists()
	case "file-contains":
		return c.checkFileContains()
	case "output-contains":
		return c.checkOutputContains()
	default:
		util.PrintFailed()
		util.PrintFix(fmt.Sprintf("Check config file and update type '%s' to a valid type.", c.Type))
		fmt.Printf("Unknown check type: '%s'\n", c.Type)
		return nil
	}
}

// checkCommandInPath checks that the configured command is in the user's path
func (c *Check) checkCommandInPath() error {
	path, err := exec.LookPath(c.Command)

	if err == nil {
		util.PrintPassed()
		if c.Verbose {
			fmt.Printf("Found %s at: %s\n", c.Command, path)
		}
		return nil
	}

	if err != nil {
		util.PrintFailed()
		util.PrintFix(c.Fix)
		if c.Verbose {
			fmt.Print(err)
		}
	}

	return nil
}

// checkFileExists checks that the configured file exists on the user's system
func (c *Check) checkFileExists() error {
	cleanPath, err := util.CleanHome(c.Path)
	if err != nil {
		util.PrintFailed()
		return err
	}

	_, err = os.Stat(cleanPath)
	if err == nil {
		util.PrintPassed()
		if c.Verbose {
			fmt.Printf("Found %s!\n", cleanPath)
		}
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		util.PrintFailed()
		util.PrintFix(c.Fix)
		if c.Verbose {
			fmt.Println(err)
		}
		return nil
	}

	util.PrintFailed()
	return err
}

// checkFileContains checks that the configured file contains specific text
func (c *Check) checkFileContains() error {
	cleanPath, err := util.CleanHome(c.Path)
	if err != nil {
		util.PrintFailed()
		return err
	}

	f, err := os.ReadFile(cleanPath)

	if err != nil {
		util.PrintFailed()
		return err
	}

	found := strings.Contains(string(f), c.Content)

	if found {
		util.PrintPassed()
		if c.Verbose {
			fmt.Printf("Found '%s' in %s!\n", c.Content, cleanPath)
		}
		return nil
	}

	if !found {
		util.PrintFailed()
		util.PrintFix(c.Fix)
		if c.Verbose {
			fmt.Printf("Could not find '%s' in %s!\n", c.Content, cleanPath)
		}
	}

	return nil
}

// checkFileContains checks that the configured command outputs specific text
func (c *Check) checkOutputContains() error {
	wrappedCmd := []string{"bash", "-c", c.Command}

	cmd := exec.Command(wrappedCmd[0], wrappedCmd[1:]...)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	if err := cmd.Run(); err != nil {
		util.PrintFailed()
		return err
	}

	found := strings.Contains(cmdOutput.String(), c.Content)

	if found {
		util.PrintPassed()
		if c.Verbose {
			fmt.Printf("Found '%s' in '%s' output!\n", c.Content, c.Command)
			fmt.Print(cmdOutput.String())
		}
		return nil
	}

	if !found {
		util.PrintFailed()
		util.PrintFix(c.Fix)
		if c.Verbose {
			fmt.Printf("Could not find '%s' in '%s' output!\n", c.Content, c.Command)
			fmt.Print(cmdOutput.String())
		}
	}

	return nil
}

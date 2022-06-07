package util

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

// CleanHome takes a path and cleans up any references to `~` to point to the
// current user's home directory
func CleanHome(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return strings.Replace(path, "~", home, -1), nil
}

// PrintDescription prints a formatted description
func PrintDescription(description string) {
	c := color.New(color.FgHiBlue)
	c.Printf("\n%s... ", description)
}

// PrintFailed prints a formatted failure message
func PrintFailed() {
	c := color.New(color.FgHiRed)
	c.Printf("failed!\n")
}

// PrintPassed prints a formatted passed message
func PrintPassed() {
	c := color.New(color.FgHiGreen)
	c.Printf("passed!\n")
}

// PrintFix prints a formatted fix message
func PrintFix(text string) {
	c := color.New(color.FgHiYellow)
	c.Printf("%s\n", text)
}

package util

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/cjonesy/doctor/internal/logger"
	"github.com/fatih/color"
)

var (
	// outputMutex ensures thread-safe output
	outputMutex sync.Mutex
)

// CleanHome takes a path and cleans up any references to `~` to point to the
// current user's home directory
func CleanHome(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(path, "~", home), nil
}

// PrintDescription prints a formatted description and logs it (thread-safe)
func PrintDescription(ctx context.Context, log *logger.Logger, description string) {
	outputMutex.Lock()
	defer outputMutex.Unlock()

	c := color.New(color.FgHiBlue)
	c.Printf("\n%s... ", description)

	if log != nil {
		log.DebugContext(ctx, "running check", "description", description)
	}
}

// PrintFailed prints a formatted failure message (thread-safe)
func PrintFailed() {
	outputMutex.Lock()
	defer outputMutex.Unlock()

	c := color.New(color.FgHiRed)
	c.Printf("failed!\n")
}

// PrintPassed prints a formatted passed message (thread-safe)
func PrintPassed() {
	outputMutex.Lock()
	defer outputMutex.Unlock()

	c := color.New(color.FgHiGreen)
	c.Printf("passed!\n")
}

// PrintFix prints a formatted fix message (thread-safe)
func PrintFix(text string) {
	outputMutex.Lock()
	defer outputMutex.Unlock()

	c := color.New(color.FgHiYellow)
	c.Printf("%s\n", text)
}

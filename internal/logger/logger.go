package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger wraps slog for application-wide logging
type Logger struct {
	*slog.Logger
}

// New creates a new logger with the specified options
func New(verbose bool, jsonOutput bool) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if verbose {
		opts.Level = slog.LevelDebug
	}

	if jsonOutput {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// CheckStarted logs when a check starts
func (l *Logger) CheckStarted(ctx context.Context, checkType, description string) {
	l.InfoContext(ctx, "check started",
		"type", checkType,
		"description", description,
	)
}

// CheckPassed logs when a check passes
func (l *Logger) CheckPassed(ctx context.Context, checkType, description string) {
	l.InfoContext(ctx, "check passed",
		"type", checkType,
		"description", description,
		"result", "pass",
	)
}

// CheckFailed logs when a check fails
func (l *Logger) CheckFailed(ctx context.Context, checkType, description, fix string) {
	l.WarnContext(ctx, "check failed",
		"type", checkType,
		"description", description,
		"result", "fail",
		"fix", fix,
	)
}

// CheckError logs when a check encounters an error
func (l *Logger) CheckError(ctx context.Context, checkType, description string, err error) {
	l.ErrorContext(ctx, "check error",
		"type", checkType,
		"description", description,
		"error", err,
	)
}

package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrConfigNotFound is returned when config file cannot be located
	ErrConfigNotFound = errors.New("config file not found")

	// ErrConfigInvalid is returned when config file is malformed
	ErrConfigInvalid = errors.New("config file is invalid")

	// ErrCheckFailed is returned when a check fails (not an error, expected outcome)
	ErrCheckFailed = errors.New("check failed")
)

// CheckError wraps errors from check execution with context
type CheckError struct {
	CheckType   string
	Description string
	Err         error
}

func (e *CheckError) Error() string {
	return fmt.Sprintf("check '%s' (%s) failed: %v", e.Description, e.CheckType, e.Err)
}

func (e *CheckError) Unwrap() error {
	return e.Err
}

// NewCheckError creates a new CheckError
func NewCheckError(checkType, description string, err error) *CheckError {
	return &CheckError{
		CheckType:   checkType,
		Description: description,
		Err:         err,
	}
}

// ConfigError wraps configuration-related errors
type ConfigError struct {
	Path string
	Err  error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error at '%s': %v", e.Path, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new ConfigError
func NewConfigError(path string, err error) *ConfigError {
	return &ConfigError{
		Path: path,
		Err:  err,
	}
}

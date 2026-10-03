package config

import (
	"fmt"
	"os"
	"strings"
)

// isStandardOutput checks if output is a standard stream alias.
func isStandardOutput(norm string) bool {
	switch norm {
	case "", OutputStdout, OutputStderr, OutputDiscard:
		return true
	}
	return false
}

// checkOutputFile verifies that the target file path can be opened and closed.
func checkOutputFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return fmt.Errorf("%w: failed to open file %q: %v", ErrInvalidOutput, path, err)
	}
	if closeErr := f.Close(); closeErr != nil {
		return fmt.Errorf("%w: failed to close file %q: %v", ErrInvalidOutput, path, closeErr)
	}
	return nil
}

// validateOutput checks whether the output destination is valid.
func (c *Config) validateOutput() error {
	norm := strings.ToLower(strings.TrimSpace(c.Output))
	if isStandardOutput(norm) {
		return nil
	}
	if isLikelyFilePath(c.Output) {
		return checkOutputFile(c.Output)
	}
	return fmt.Errorf("%w: %q (supported: stdout, stderr, discard, or file path)", ErrInvalidOutput, c.Output)
}

// validateFormat checks whether the configured format is supported.
func (c *Config) validateFormat() error {
	_, err := ParseFormat(c.Format)
	return err
}

// validateLevel checks whether the configured level is valid.
func (c *Config) validateLevel() error {
	_, err := ParseLevel(c.Level)
	return err
}

func (c *Config) validateFields() error {
	if err := c.validateOutput(); err != nil {
		return err
	}
	if err := c.validateFormat(); err != nil {
		return err
	}
	return c.validateLevel()
}

// Validate validates the configuration values.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}
	return c.validateFields()
}

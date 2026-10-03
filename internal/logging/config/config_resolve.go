package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func resolveStderrOrDiscard(norm string) (io.Writer, bool) {
	if norm == OutputStderr {
		return os.Stderr, true
	}
	if norm == OutputDiscard {
		return io.Discard, true
	}
	return nil, false
}

func resolveStandardOutput(norm string) (io.Writer, bool) {
	if norm == "" || norm == OutputStdout {
		return os.Stdout, true
	}
	return resolveStderrOrDiscard(norm)
}

func openOutputFile(path string) (io.Writer, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to open file %q: %v", ErrInvalidOutput, path, err)
	}
	return f, nil
}

// ResolveOutput resolves the output destination to an io.Writer.
func (c *Config) ResolveOutput() (io.Writer, error) {
	if c == nil {
		return nil, ErrNilConfig
	}
	norm := strings.ToLower(strings.TrimSpace(c.Output))
	if out, ok := resolveStandardOutput(norm); ok {
		return out, nil
	}
	return openOutputFile(c.Output)
}

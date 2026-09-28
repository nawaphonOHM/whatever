package logging

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Format represents supported log output formats.
type Format string

// Supported output format constants.
const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// ErrInvalidFormat indicates an unsupported log format string.
var ErrInvalidFormat = errors.New("invalid log format")

// ReplaceAttrFunc specifies a function to modify attributes before logging.
type ReplaceAttrFunc = func([]string, slog.Attr) slog.Attr

// Config contains options for initializing a structured logger.
type Config struct {
	// Output is the destination for log writes. Defaults to os.Stdout if nil.
	Output io.Writer

	// ReplaceAttr allows customizing log attributes before writing.
	ReplaceAttr ReplaceAttrFunc

	// ExitFunc is the function invoked on Fatal/FatalContext. Defaults to os.Exit.
	ExitFunc func(int)

	// Handler allows supplying a pre-configured slog.Handler directly.
	Handler slog.Handler

	// Level specifies the minimum severity level to log (TRACE, DEBUG, INFO, WARN, ERROR, FATAL).
	Level Level

	// Format specifies the output format (json or text). Defaults to FormatJSON.
	Format Format

	// AddSource attaches caller file and line numbers to records when true.
	AddSource bool

	// DisableTraceCorrelation disables automatic OpenTelemetry trace_id and span_id enrichment.
	DisableTraceCorrelation bool
}

// DefaultConfig returns recommended logger configuration defaults.
func DefaultConfig() Config {
	return Config{
		Output:                  os.Stdout,
		ExitFunc:                os.Exit,
		Level:                   LevelInfo,
		Format:                  FormatJSON,
		AddSource:               false,
		DisableTraceCorrelation: false,
	}
}

// isSupportedFormat checks if format string is empty or one of supported formats.
func isSupportedFormat(f Format) bool {
	norm := strings.ToLower(strings.TrimSpace(string(f)))
	return norm == "" || norm == string(FormatJSON) || norm == string(FormatText)
}

// validateFormat checks whether the configured format is supported.
func (c *Config) validateFormat() error {
	if isSupportedFormat(c.Format) {
		return nil
	}
	return fmt.Errorf("%w: %q (supported: %s, %s)", ErrInvalidFormat, c.Format, FormatJSON, FormatText)
}

// validateLevel checks whether the configured level is valid.
func (c *Config) validateLevel() error {
	if c.Level == "" {
		return nil
	}
	_, err := ParseLevel(string(c.Level))
	return err
}

// Validate validates the configuration values.
func (c *Config) Validate() error {
	if err := c.validateFormat(); err != nil {
		return err
	}
	return c.validateLevel()
}

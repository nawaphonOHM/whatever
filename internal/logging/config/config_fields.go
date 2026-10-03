package config

import (
	"errors"
)

// Supported output format constants.
const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Format represents supported log output formats.
type Format string

// Supported standard output stream destination constants.
const (
	OutputStdout  = "stdout"
	OutputStderr  = "stderr"
	OutputDiscard = "discard"
)

// Default logging configuration constants.
const (
	DefaultOutput                  = "stdout"
	DefaultLevel                   = "info"
	DefaultFormat                  = "text"
	DefaultAddSource               = false
	DefaultDisableTraceCorrelation = false
)

// Sentinel errors for logging configuration validation.
var (
	// ErrNilConfig is returned when validating a nil configuration.
	ErrNilConfig = errors.New("logging config cannot be nil")

	// ErrInvalidFormat is returned when the log format is unsupported.
	ErrInvalidFormat = errors.New("invalid log format")

	// ErrInvalidLevel is returned when the log level is unrecognized.
	ErrInvalidLevel = errors.New("invalid log level")

	// ErrInvalidOutput is returned when the log output destination is unsupported or invalid.
	ErrInvalidOutput = errors.New("invalid log output destination")
)

// Config defines the configuration options for application logging.
// All environment variable keys use the OHM9996_LOGGING_ prefix.
type Config struct {
	Output                  string `env:"OHM9996_LOGGING_OUTPUT" envDefault:"stdout"`
	Level                   string `env:"OHM9996_LOGGING_LEVEL" envDefault:"info"`
	Format                  string `env:"OHM9996_LOGGING_FORMAT" envDefault:"text"`
	AddSource               bool   `env:"OHM9996_LOGGING_ADD_SOURCE" envDefault:"false"`
	DisableTraceCorrelation bool   `env:"OHM9996_LOGGING_DISABLE_TRACE_CORRELATION" envDefault:"false"`
}

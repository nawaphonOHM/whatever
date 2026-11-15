package azurite

import (
	"github.com/testcontainers/testcontainers-go"

	internalazurite "github.com/nawaphonOHM/whatever/v3/internal/testcontainers/azurite"
)

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return internalazurite.DefaultOptions()
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	return internalazurite.NewOptions(opts...)
}

// WithImage configures the Docker image used for the Azurite container.
func WithImage(image string) Option {
	return internalazurite.WithImage(image)
}

// WithEnv sets a custom environment variable for the Azurite container.
func WithEnv(key, value string) Option {
	return internalazurite.WithEnv(key, value)
}

// WithContainerOptions appends raw Testcontainers customizers.
func WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option {
	return internalazurite.WithContainerOptions(opts...)
}

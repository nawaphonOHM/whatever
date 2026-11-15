package azurite

import "github.com/testcontainers/testcontainers-go"

// DefaultImage is the default image used for the Azurite emulator.
const DefaultImage = "mcr.microsoft.com/azure-storage/azurite:3.37.0"

// Option defines a functional option for configuring an Azurite container.
type Option func(*Options)

// Options holds configuration settings for running an Azurite container.
type Options struct {
	Image            string
	Env              map[string]string
	ContainerOptions []testcontainers.ContainerCustomizer
}

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return &Options{
		Image: DefaultImage,
		Env:   make(map[string]string),
	}
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	o := DefaultOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}

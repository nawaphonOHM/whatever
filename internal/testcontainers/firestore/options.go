package firestore

import (
	"github.com/testcontainers/testcontainers-go"
)

// Option defines a functional option for configuring a Firestore container.
type Option func(*Options)

// Options holds configuration settings for running a Firestore container.
type Options struct {
	Image            string
	ProjectID        string
	Env              map[string]string
	ContainerOptions []testcontainers.ContainerCustomizer
	DatastoreMode    bool
}

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return &Options{
		Image:     DefaultImage,
		ProjectID: DefaultProjectID,
		Env:       make(map[string]string),
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

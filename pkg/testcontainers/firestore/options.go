package firestore

import (
	"github.com/testcontainers/testcontainers-go"

	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/firestore"
)

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return firestore.DefaultOptions()
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	return firestore.NewOptions(opts...)
}

// WithImage configures the Docker image used for the Firestore container.
func WithImage(image string) Option {
	return firestore.WithImage(image)
}

// WithProjectID sets the Google Cloud project ID for the Firestore container.
func WithProjectID(projectID string) Option {
	return firestore.WithProjectID(projectID)
}

// WithDatastoreMode enables or disables Firestore datastore mode.
func WithDatastoreMode(enabled ...bool) Option {
	return firestore.WithDatastoreMode(enabled...)
}

// WithEnv sets a custom environment variable for the Firestore container.
func WithEnv(key, value string) Option {
	return firestore.WithEnv(key, value)
}

// WithContainerOptions appends raw Testcontainers customizers.
func WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option {
	return firestore.WithContainerOptions(opts...)
}

package mongodb

import (
	"github.com/testcontainers/testcontainers-go"

	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/mongodb"
)

// DefaultOptions returns an Options struct populated with default settings.
func DefaultOptions() *Options {
	return mongodb.DefaultOptions()
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	return mongodb.NewOptions(opts...)
}

// WithImage configures the Docker image used for the MongoDB container.
func WithImage(image string) Option {
	return mongodb.WithImage(image)
}

// WithUsername sets the root username for MongoDB authentication.
func WithUsername(username string) Option {
	return mongodb.WithUsername(username)
}

// WithPassword sets the root password for MongoDB authentication.
func WithPassword(password string) Option {
	return mongodb.WithPassword(password)
}

// WithDatabase sets the default database name created on container initialization.
func WithDatabase(database string) Option {
	return mongodb.WithDatabase(database)
}

// WithReplicaSet configures the single-node replica set name.
func WithReplicaSet(replicaSet string) Option {
	return mongodb.WithReplicaSet(replicaSet)
}

// WithEnv sets a custom environment variable for the MongoDB container.
func WithEnv(key, value string) Option {
	return mongodb.WithEnv(key, value)
}

// WithContainerOptions appends raw Testcontainers customizers.
func WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option {
	return mongodb.WithContainerOptions(opts...)
}

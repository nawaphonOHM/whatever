package mongodb

import (
	"github.com/testcontainers/testcontainers-go"
)

// WithImage configures the Docker image used for the MongoDB container.
func WithImage(image string) Option {
	return func(o *Options) {
		o.Image = image
	}
}

// WithUsername sets the root username for MongoDB authentication.
func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

// WithPassword sets the root password for MongoDB authentication.
func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// WithDatabase sets the default database name created on container initialization.
func WithDatabase(database string) Option {
	return func(o *Options) {
		o.Database = database
	}
}

// WithReplicaSet configures the single-node replica set name.
func WithReplicaSet(replicaSet string) Option {
	return func(o *Options) {
		o.ReplicaSet = replicaSet
	}
}

// WithEnv sets a custom environment variable for the MongoDB container.
func WithEnv(key, value string) Option {
	return func(o *Options) {
		if o.Env == nil {
			o.Env = make(map[string]string)
		}
		o.Env[key] = value
	}
}

// WithContainerOptions appends raw Testcontainers customizers.
func WithContainerOptions(opts ...testcontainers.ContainerCustomizer) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.ContainerOptions = append(o.ContainerOptions, opt)
			}
		}
	}
}

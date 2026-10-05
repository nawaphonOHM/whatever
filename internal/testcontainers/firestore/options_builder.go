package firestore

import (
	"github.com/testcontainers/testcontainers-go"
)

// WithImage configures the Docker image used for the Firestore container.
func WithImage(image string) Option {
	return func(o *Options) {
		o.Image = image
	}
}

// WithProjectID sets the Google Cloud project ID for the Firestore container.
func WithProjectID(projectID string) Option {
	return func(o *Options) {
		o.ProjectID = projectID
	}
}

// WithDatastoreMode enables or disables Firestore datastore mode.
func WithDatastoreMode(enabled ...bool) Option {
	return func(o *Options) {
		if len(enabled) == 0 {
			o.DatastoreMode = true
			return
		}
		o.DatastoreMode = enabled[0]
	}
}

// WithEnv sets a custom environment variable for the Firestore container.
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

package azurite

import "github.com/testcontainers/testcontainers-go"

// WithImage configures the Docker image used for the Azurite container.
func WithImage(image string) Option {
	return func(o *Options) {
		o.Image = image
	}
}

// WithEnv sets a custom environment variable for the Azurite container.
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

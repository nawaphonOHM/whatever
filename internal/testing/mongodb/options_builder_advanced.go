package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// WithConnectTimeout sets the timeout for initial connection establishment.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ConnectTimeout = timeout
	}
}

// WithServerSelectionTimeout sets the timeout for server selection.
func WithServerSelectionTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ServerSelectionTimeout = timeout
	}
}

// WithSocketTimeout sets the timeout for socket read/write operations.
func WithSocketTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.SocketTimeout = timeout
	}
}

// WithMaxConnIdleTime sets the maximum connection idle duration.
func WithMaxConnIdleTime(idleTime time.Duration) Option {
	return func(o *Options) {
		o.MaxConnIdleTime = idleTime
	}
}

// WithPoolLimits sets the minimum and maximum connection pool sizes.
func WithPoolLimits(minPoolSize, maxPoolSize uint64) Option {
	return func(o *Options) {
		o.MinPoolSize = minPoolSize
		o.MaxPoolSize = maxPoolSize
	}
}

// WithTLS sets whether TLS encryption is enabled.
func WithTLS(enableTLS bool) Option {
	return func(o *Options) {
		o.EnableTLS = enableTLS
	}
}

// WithDirectConnection sets whether to force a direct connection to a single host.
func WithDirectConnection(direct bool) Option {
	return func(o *Options) {
		o.DirectConnection = direct
	}
}

// WithUUIDRepresentation sets the binary UUID encoding representation.
func WithUUIDRepresentation(repr string) Option {
	return func(o *Options) {
		o.UUIDRepresentation = repr
	}
}

// WithDriverOptions appends raw mongo-driver ClientOptions.
func WithDriverOptions(opts ...*options.ClientOptions) Option {
	return func(o *Options) {
		for _, opt := range opts {
			if opt != nil {
				o.DriverOptions = append(o.DriverOptions, opt)
			}
		}
	}
}

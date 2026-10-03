package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/testing/mongodb"
)

// WithConnectTimeout sets the timeout for initial connection establishment.
func WithConnectTimeout(timeout time.Duration) Option {
	return mongodb.WithConnectTimeout(timeout)
}

// WithServerSelectionTimeout sets the timeout for server selection.
func WithServerSelectionTimeout(timeout time.Duration) Option {
	return mongodb.WithServerSelectionTimeout(timeout)
}

// WithSocketTimeout sets the timeout for socket read/write operations.
func WithSocketTimeout(timeout time.Duration) Option {
	return mongodb.WithSocketTimeout(timeout)
}

// WithMaxConnIdleTime sets the maximum connection idle duration.
func WithMaxConnIdleTime(idleTime time.Duration) Option {
	return mongodb.WithMaxConnIdleTime(idleTime)
}

// WithPoolLimits sets the minimum and maximum connection pool sizes.
func WithPoolLimits(minPoolSize, maxPoolSize uint64) Option {
	return mongodb.WithPoolLimits(minPoolSize, maxPoolSize)
}

// WithPing sets whether to verify connectivity with a ping on connect.
func WithPing(enablePing bool) Option {
	return mongodb.WithPing(enablePing)
}

// WithTLS sets whether TLS encryption is enabled.
func WithTLS(enableTLS bool) Option {
	return mongodb.WithTLS(enableTLS)
}

// WithDirectConnection sets whether to force a direct connection to a single host.
func WithDirectConnection(direct bool) Option {
	return mongodb.WithDirectConnection(direct)
}

// WithUUIDRepresentation sets the binary UUID encoding representation.
func WithUUIDRepresentation(repr string) Option {
	return mongodb.WithUUIDRepresentation(repr)
}

// WithDriverOptions appends raw mongo-driver ClientOptions.
func WithDriverOptions(opts ...*options.ClientOptions) Option {
	return mongodb.WithDriverOptions(opts...)
}

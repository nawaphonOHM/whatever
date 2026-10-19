package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/v3/internal/testing/mongodb"
)

// Option defines a functional option for configuring MongoDB test connections.
type Option = mongodb.Option

// Options holds configuration settings for test MongoDB connections.
type Options = mongodb.Options

// DefaultOptions returns an Options struct populated with recommended test defaults.
func DefaultOptions() *Options {
	return mongodb.DefaultOptions()
}

// NewOptions creates an Options struct by evaluating functional options over defaults.
func NewOptions(opts ...Option) *Options {
	return mongodb.NewOptions(opts...)
}

// WithURI sets the direct connection URI string.
func WithURI(uri string) Option {
	return mongodb.WithURI(uri)
}

// WithHost sets the MongoDB host address.
func WithHost(host string) Option {
	return mongodb.WithHost(host)
}

// WithPort sets the MongoDB connection port.
func WithPort(port int) Option {
	return mongodb.WithPort(port)
}

// WithProtocol sets the connection scheme (e.g., mongodb or mongodb+srv).
func WithProtocol(protocol string) Option {
	return mongodb.WithProtocol(protocol)
}

// WithDatabase sets the default database name for the test client.
func WithDatabase(database string) Option {
	return mongodb.WithDatabase(database)
}

// WithUsername sets the authentication username.
func WithUsername(username string) Option {
	return mongodb.WithUsername(username)
}

// WithPassword sets the authentication password.
func WithPassword(password string) Option {
	return mongodb.WithPassword(password)
}

// WithAuthSource sets the authentication source database.
func WithAuthSource(authSource string) Option {
	return mongodb.WithAuthSource(authSource)
}

// WithAppName sets the client application name.
func WithAppName(appName string) Option {
	return mongodb.WithAppName(appName)
}

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

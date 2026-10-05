package mongodb

import (
	"github.com/nawaphonOHM/whatever/v2/internal/testing/mongodb"
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

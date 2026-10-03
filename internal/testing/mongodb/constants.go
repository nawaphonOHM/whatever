// Package mongodb provides internal testing utilities, fixtures, options,
// and connection helpers for MongoDB test suites.
package mongodb

import "time"

const (
	// DefaultHost is the default MongoDB host for test connections.
	DefaultHost = "localhost"
	// DefaultPort is the default MongoDB port for test connections.
	DefaultPort = 27017
	// DefaultProtocol is the default MongoDB connection scheme.
	DefaultProtocol = "mongodb"
	// DefaultConnectTimeout is the default connection timeout for test clients.
	DefaultConnectTimeout = 10 * time.Second
	// DefaultServerSelectionTimeout is the default server selection timeout for test clients.
	DefaultServerSelectionTimeout = 5 * time.Second
	// DefaultSocketTimeout is the default socket/operation timeout for test clients.
	DefaultSocketTimeout = 10 * time.Second
	// DefaultMaxConnIdleTime is the default maximum connection idle time for test clients.
	DefaultMaxConnIdleTime = 10 * time.Minute
	// DefaultMaxPoolSize is the default maximum connection pool capacity.
	DefaultMaxPoolSize uint64 = 100
	// DefaultMinPoolSize is the default minimum connection pool capacity.
	DefaultMinPoolSize uint64 = 5
	// DefaultUUIDRepresentation is the default UUID encoding representation.
	DefaultUUIDRepresentation = "unspecified"
)

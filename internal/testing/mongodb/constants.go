// Package mongodb provides internal testing utilities, fixtures, options,
// and connection helpers for MongoDB test suites.
package mongodb

import (
	"errors"
	"time"
)

const (
	// DefaultHost is the default MongoDB host for test connections.
	DefaultHost = "localhost"
	// DefaultPort is the default MongoDB port for test connections.
	DefaultPort = 27017
	// DefaultFirestorePort is the default port for Firestore MongoDB endpoints.
	DefaultFirestorePort = 443
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

var (
	// ErrNilClient is returned when an operation is performed on a nil TestClient.
	ErrNilClient = errors.New("mongodb test client is nil")
	// ErrNilConfig is returned when an operation receives nil options or configuration.
	ErrNilConfig = errors.New("mongodb config cannot be nil")
	// ErrEmptyURI is returned when attempting to connect with an empty URI string.
	ErrEmptyURI = errors.New("mongodb uri cannot be empty")
	// ErrInvalidPort is returned when the configured port is outside the valid range.
	ErrInvalidPort = errors.New("mongodb port must be between 1 and 65535")
	// ErrNilTestingTB is returned when a testing helper receives a nil testing.TB instance.
	ErrNilTestingTB = errors.New("testing.TB cannot be nil")
)

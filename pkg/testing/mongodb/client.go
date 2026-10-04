// Package mongodb provides public MongoDB testing utilities, type aliases,
// functional options, and testing.TB lifecycle integration for automated test suites.
package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/testing/mongodb"
)

// TestClient wraps mongo client handles with testing and fixture utilities.
type TestClient = mongodb.TestClient

// Option defines a functional option for configuring MongoDB test connections.
type Option = mongodb.Option

// Options holds configuration settings for test MongoDB connections.
type Options = mongodb.Options

// Constants for MongoDB testing default settings.
const (
	// DefaultHost is the default MongoDB host for test connections.
	DefaultHost = mongodb.DefaultHost
	// DefaultPort is the default MongoDB port for test connections.
	DefaultPort = mongodb.DefaultPort
	// DefaultProtocol is the default MongoDB connection scheme.
	DefaultProtocol = mongodb.DefaultProtocol
	// DefaultConnectTimeout is the default connection timeout for test clients.
	DefaultConnectTimeout = mongodb.DefaultConnectTimeout
	// DefaultServerSelectionTimeout is the default server selection timeout for test clients.
	DefaultServerSelectionTimeout = mongodb.DefaultServerSelectionTimeout
	// DefaultSocketTimeout is the default socket/operation timeout for test clients.
	DefaultSocketTimeout = mongodb.DefaultSocketTimeout
	// DefaultMaxConnIdleTime is the default maximum connection idle time for test clients.
	DefaultMaxConnIdleTime = mongodb.DefaultMaxConnIdleTime
	// DefaultMaxPoolSize is the default maximum connection pool capacity.
	DefaultMaxPoolSize = mongodb.DefaultMaxPoolSize
	// DefaultMinPoolSize is the default minimum connection pool capacity.
	DefaultMinPoolSize = mongodb.DefaultMinPoolSize
	// DefaultUUIDRepresentation is the default UUID encoding representation.
	DefaultUUIDRepresentation = mongodb.DefaultUUIDRepresentation
)

// Sentinel errors for MongoDB testing client operations.
var (
	// ErrNilClient is returned when an operation is performed on a nil TestClient.
	ErrNilClient = mongodb.ErrNilClient
	// ErrNilConfig is returned when an operation receives nil options or configuration.
	ErrNilConfig = mongodb.ErrNilConfig
	// ErrEmptyURI is returned when attempting to connect with an empty URI string.
	ErrEmptyURI = mongodb.ErrEmptyURI
	// ErrInvalidPort is returned when the configured port is outside the valid range.
	ErrInvalidPort = mongodb.ErrInvalidPort
	// ErrNilTestingTB is returned when a testing helper receives a nil testing.TB instance.
	ErrNilTestingTB = mongodb.ErrNilTestingTB
)

// NewClient creates a new TestClient wrapping an official mongo.Client.
func NewClient(rawClient *mongo.Client, defaultDatabase string) *TestClient {
	return mongodb.NewClient(rawClient, defaultDatabase)
}

// Package mongodb provides the public entrypoint and managed client for
// MongoDB connections.
package mongodb

import (
	"context"

	intmongo "github.com/nawaphonOHM/whatever/internal/mongodb"
)

// Client wraps the official mongo.Client and provides managed database and
// collection access.
type Client = intmongo.Client

// Sentinel errors for MongoDB client operations.
var (
	ErrNilClient = intmongo.ErrNilClient
	ErrNilConfig = intmongo.ErrNilConfig
)

// Connect loads MongoDB configuration from environment variables with the
// OHM9996_MONGODB_ prefix, connects to the MongoDB deployment, verifies
// connectivity via Ping, and returns a managed Client.
func Connect(ctx context.Context) (*Client, error) {
	return intmongo.Connect(ctx)
}

// Package mongodb provides the public entrypoint and managed client for
// MongoDB connections.
package mongodb

import (
	"context"

	"github.com/nawaphonOHM/whatever/internal/mongodb/client"
	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

// Client wraps the official mongo.Client and provides managed database and
// collection access.
type Client = client.Client

// Sentinel errors for MongoDB client operations.
var (
	ErrNilClient = client.ErrNilClient
	ErrNilConfig = config.ErrNilConfig
)

// Connect loads MongoDB configuration from environment variables with the
// OHM9996_MONGODB_ prefix, connects to the MongoDB deployment, verifies
// connectivity via Ping, and returns a managed Client.
func Connect(ctx context.Context) (*Client, error) {
	return client.Connect(ctx)
}

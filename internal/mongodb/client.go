package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Sentinel errors for MongoDB client operations.
var (
	ErrNilClient = errors.New("mongodb client is not initialized")
	ErrNilConfig = errors.New("mongodb config cannot be nil")
)

// Client wraps the official mongo.Client and manages connection pooling,
// lifecycle operations, and default database resolution.
type Client struct {
	rawClient       *mongo.Client
	defaultDatabase string
}

// NewClient creates a new Client wrapping an existing mongo.Client with a
// default database.
func NewClient(rawClient *mongo.Client, defaultDatabase string) *Client {
	return &Client{
		rawClient:       rawClient,
		defaultDatabase: defaultDatabase,
	}
}

// resolveDBName returns the target database name or default fallback.
func (c *Client) resolveDBName(name ...string) string {
	if len(name) > 0 && name[0] != "" {
		return name[0]
	}
	return c.defaultDatabase
}

// Database returns a handle to the specified database.
// If no database name or an empty name is provided, it falls back to the
// configured default database. Returns nil if the client is not initialized.
func (c *Client) Database(name ...string) *mongo.Database {
	if c == nil || c.rawClient == nil {
		return nil
	}
	return c.rawClient.Database(c.resolveDBName(name...))
}

// Collection returns a handle for a collection in the specified database.
// If dbName is omitted or empty, it falls back to default database.
func (c *Client) Collection(name string, dbName ...string) *mongo.Collection {
	db := c.Database(dbName...)
	if db == nil {
		return nil
	}
	return db.Collection(name)
}

// Ping sends a ping command to verify connectivity to the MongoDB deployment.
func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.rawClient == nil {
		return ErrNilClient
	}
	return c.rawClient.Ping(ctx, readpref.Primary())
}

// Disconnect gracefully closes all sockets in the client connection pool.
func (c *Client) Disconnect(ctx context.Context) error {
	if c == nil || c.rawClient == nil {
		return ErrNilClient
	}
	return c.rawClient.Disconnect(ctx)
}

// RawClient returns the underlying official *mongo.Client handle.
func (c *Client) RawClient() *mongo.Client {
	if c == nil {
		return nil
	}
	return c.rawClient
}

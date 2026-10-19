// Package client provides the managed MongoDB client wrapper, connection
// lifecycle management, and official mongo-driver options construction.
package client

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// ErrNilClient Sentinel errors for MongoDB client operations.
var (
	ErrNilClient = errors.New("mongodb client is not initialized")
)

// Client wraps the official mongo.Client and manages connection pooling,
// lifecycle operations, and default database resolution.
type Client struct {
	rawClient       *mongo.Client
	defaultDatabase string
	isFirestore     bool
}

// NewClient creates a new Client wrapping an existing mongo.Client with a
// default database and optional Firestore mode.
func NewClient(rawClient *mongo.Client, defaultDatabase string, isFirestore ...bool) *Client {
	var firestore bool
	if len(isFirestore) > 0 {
		firestore = isFirestore[0]
	}
	return &Client{
		rawClient:       rawClient,
		defaultDatabase: defaultDatabase,
		isFirestore:     firestore,
	}
}

// IsFirestore reports whether the client is configured to connect to Google Cloud Firestore.
func (c *Client) IsFirestore() bool {
	if c == nil {
		return false
	}
	return c.isFirestore
}

func (c *Client) isNil() bool {
	return c == nil || c.rawClient == nil
}

// Ping sends a ping command to verify connectivity to the MongoDB deployment,
// or performs a Firestore dummy read query if connected to Google Cloud Firestore.
func (c *Client) Ping(ctx context.Context) error {
	if c.isNil() {
		return ErrNilClient
	}
	if c.isFirestore {
		return pingFirestore(ctx, c.rawClient, c.defaultDatabase)
	}
	return c.rawClient.Ping(ctx, readpref.Primary())
}

// Disconnect gracefully closes all sockets in the client connection pool.
func (c *Client) Disconnect(ctx context.Context) error {
	if c.isNil() {
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
	if c.isNil() {
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

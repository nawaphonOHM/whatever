package mongodb

import (
	"context"

	"github.com/nawaphonOHM/whatever/internal/mongodb/client"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// TestClient wraps mongo client handles with testing and fixture utilities.
type TestClient struct {
	client          *client.Client
	rawClient       *mongo.Client
	defaultDatabase string
}

// NewClient creates a new TestClient wrapping an official mongo.Client.
func NewClient(rawClient *mongo.Client, defaultDatabase string) *TestClient {
	return &TestClient{
		client:          client.NewClient(rawClient, defaultDatabase),
		rawClient:       rawClient,
		defaultDatabase: defaultDatabase,
	}
}

// isNil reports whether the TestClient receiver or raw client handle is nil.
func (c *TestClient) isNil() bool {
	return c == nil || c.rawClient == nil
}

// resolveDBName returns the target database name or default fallback.
func (c *TestClient) resolveDBName(name ...string) string {
	if len(name) > 0 && name[0] != "" {
		return name[0]
	}
	return c.defaultDatabase
}

// Database returns a handle to the specified database or default fallback.
func (c *TestClient) Database(name ...string) *mongo.Database {
	if c.isNil() {
		return nil
	}
	return c.rawClient.Database(c.resolveDBName(name...))
}

// Collection returns a handle for a collection in the specified database.
func (c *TestClient) Collection(name string, dbName ...string) *mongo.Collection {
	db := c.Database(dbName...)
	if db == nil {
		return nil
	}
	return db.Collection(name)
}

// Ping verifies connectivity to the MongoDB deployment.
func (c *TestClient) Ping(ctx context.Context) error {
	if c.isNil() {
		return ErrNilClient
	}
	return c.rawClient.Ping(ctx, readpref.Primary())
}

// Disconnect gracefully closes all sockets in the client connection pool.
func (c *TestClient) Disconnect(ctx context.Context) error {
	if c.isNil() {
		return ErrNilClient
	}
	return c.rawClient.Disconnect(ctx)
}

// resolveCloseContext selects explicit or default timeout context for close.
func resolveCloseContext(ctx ...context.Context) (context.Context, context.CancelFunc) {
	if len(ctx) > 0 && ctx[0] != nil {
		return ctx[0], func() {}
	}
	return context.WithTimeout(context.Background(), DefaultConnectTimeout)
}

// Close gracefully closes the client, using provided context or a default timeout.
func (c *TestClient) Close(ctx ...context.Context) error {
	if c.isNil() {
		return ErrNilClient
	}
	closeCtx, cancel := resolveCloseContext(ctx...)
	defer cancel()
	return c.rawClient.Disconnect(closeCtx)
}

// RawClient returns the underlying official *mongo.Client handle.
func (c *TestClient) RawClient() *mongo.Client {
	if c == nil {
		return nil
	}
	return c.rawClient
}

// Client returns the underlying internal *client.Client handle.
func (c *TestClient) Client() *client.Client {
	if c == nil {
		return nil
	}
	return c.client
}

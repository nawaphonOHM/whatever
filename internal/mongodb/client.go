package mongodb

import (
	"context"
	"errors"
	"fmt"

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

// NewClient creates a new Client wrapping an existing mongo.Client with a default database.
func NewClient(rawClient *mongo.Client, defaultDatabase string) *Client {
	return &Client{
		rawClient:       rawClient,
		defaultDatabase: defaultDatabase,
	}
}

// Database returns a handle to the specified database.
// If no database name or an empty name is provided, it falls back to the configured default database.
// Returns nil if the client is not initialized.
func (c *Client) Database(name ...string) *mongo.Database {
	if c == nil || c.rawClient == nil {
		return nil
	}
	dbName := c.defaultDatabase
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}
	return c.rawClient.Database(dbName)
}

// Collection returns a handle for a collection in the specified database.
// If dbName is omitted or empty, it falls back to the configured default database.
// Returns nil if the client is not initialized.
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
// This serves as an escape hatch for advanced operations such as transactions, sessions, or change streams.
func (c *Client) RawClient() *mongo.Client {
	if c == nil {
		return nil
	}
	return c.rawClient
}

// Connect loads MongoDB configuration from environment variables (OHM9969_MONGODB_*),
// establishes a connection to the MongoDB deployment, verifies connectivity with Ping,
// and returns a managed Client.
func Connect(ctx context.Context, opts ...Option) (*Client, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return ConnectWithConfig(ctx, cfg, opts...)
}

// ConnectWithConfig connects to MongoDB using the provided Config and Options.
// It verifies connectivity via Ping using the provided context.
func ConnectWithConfig(ctx context.Context, cfg *Config, opts ...Option) (*Client, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mongodb config: %w", err)
	}

	optionsContainer := NewOptions(opts...)
	clientOptions := BuildClientOptions(cfg, optionsContainer.DriverOptions...)

	rawClient, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	client := NewClient(rawClient, cfg.Database)
	if err := client.Ping(ctx); err != nil {
		_ = rawClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}

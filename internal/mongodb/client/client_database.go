package client

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

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

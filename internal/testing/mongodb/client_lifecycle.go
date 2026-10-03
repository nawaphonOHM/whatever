package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// systemCollectionPrefix defines the prefix identifying internal system collections.
const systemCollectionPrefix = "system."

// isSystemCollection reports whether a collection name indicates an internal system collection.
func isSystemCollection(name string) bool {
	return strings.HasPrefix(name, systemCollectionPrefix)
}

// filterNonSystemCollections filters out system collections from collection names.
func filterNonSystemCollections(names []string) []string {
	result := make([]string, 0, len(names))
	for _, name := range names {
		if !isSystemCollection(name) {
			result = append(result, name)
		}
	}
	return result
}

// DropDatabase drops the specified database or default fallback database.
func (c *TestClient) DropDatabase(ctx context.Context, name ...string) error {
	if c.isNil() {
		return ErrNilClient
	}
	return c.Database(name...).Drop(ctx)
}

// DropCollection drops the specified collection from target or default database.
func (c *TestClient) DropCollection(ctx context.Context, collection string, dbName ...string) error {
	if c.isNil() {
		return ErrNilClient
	}
	return c.Collection(collection, dbName...).Drop(ctx)
}

// rawCollectionNames fetches raw collection names from the resolved database.
func (c *TestClient) rawCollectionNames(ctx context.Context, dbName ...string) ([]string, error) {
	db := c.Database(dbName...)
	return db.ListCollectionNames(ctx, bson.D{})
}

// ListCollectionNames returns all non-system collection names in the target database.
func (c *TestClient) ListCollectionNames(ctx context.Context, dbName ...string) ([]string, error) {
	if c.isNil() {
		return nil, ErrNilClient
	}
	rawNames, err := c.rawCollectionNames(ctx, dbName...)
	if err != nil {
		return nil, err
	}
	return filterNonSystemCollections(rawNames), nil
}

// truncateNamed removes all documents from a single named collection.
func (c *TestClient) truncateNamed(ctx context.Context, name string) error {
	if name == "" {
		return nil
	}
	coll := c.Collection(name)
	if coll == nil {
		return ErrNilClient
	}
	_, err := coll.DeleteMany(ctx, bson.D{})
	return err
}

// resolveTruncateTargets returns target collections to truncate.
func (c *TestClient) resolveTruncateTargets(ctx context.Context, collections []string) ([]string, error) {
	if len(collections) > 0 {
		return collections, nil
	}
	return c.ListCollectionNames(ctx)
}

// truncateAll deletes all documents across the provided collection names.
func (c *TestClient) truncateAll(ctx context.Context, targets []string) error {
	for _, name := range targets {
		if err := c.truncateNamed(ctx, name); err != nil {
			return err
		}
	}
	return nil
}

// TruncateCollections removes all documents from specified collections or all
// non-system collections in the default database.
func (c *TestClient) TruncateCollections(ctx context.Context, collections ...string) error {
	if c.isNil() {
		return ErrNilClient
	}
	targets, err := c.resolveTruncateTargets(ctx, collections)
	if err != nil {
		return err
	}
	return c.truncateAll(ctx, targets)
}

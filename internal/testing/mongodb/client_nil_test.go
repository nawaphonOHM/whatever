package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNilDB   = "test_nil_db"
	testNilColl = "test_nil_coll"
)

// verifyNilAccessors checks accessor methods on nil or uninitialized clients.
func verifyNilAccessors(t *testing.T, c *TestClient) {
	assert.Nil(t, c.Database())
	assert.Nil(t, c.Database(testNilDB))
	assert.Nil(t, c.Collection(testNilColl))
	assert.Nil(t, c.Collection(testNilColl, testNilDB))
	assert.Nil(t, c.RawClient())
	assert.Nil(t, c.Client())
}

// verifyNilLifecycle checks connection lifecycle methods on nil or uninitialized clients.
func verifyNilLifecycle(t *testing.T, c *TestClient) {
	ctx := context.Background()

	assert.ErrorIs(t, c.Ping(ctx), ErrNilClient)
	assert.ErrorIs(t, c.Disconnect(ctx), ErrNilClient)
	assert.ErrorIs(t, c.Close(), ErrNilClient)
	assert.ErrorIs(t, c.Close(ctx), ErrNilClient)
}

// verifyNilFixtures checks fixture management methods on nil or uninitialized clients.
func verifyNilFixtures(t *testing.T, c *TestClient) {
	ctx := context.Background()

	assert.ErrorIs(t, c.DropDatabase(ctx), ErrNilClient)
	assert.ErrorIs(t, c.DropDatabase(ctx, testNilDB), ErrNilClient)
	assert.ErrorIs(t, c.DropCollection(ctx, testNilColl), ErrNilClient)
	assert.ErrorIs(t, c.DropCollection(ctx, testNilColl, testNilDB), ErrNilClient)

	names, err := c.ListCollectionNames(ctx)
	assert.Nil(t, names)
	assert.ErrorIs(t, err, ErrNilClient)

	assert.ErrorIs(t, c.TruncateCollections(ctx), ErrNilClient)
	assert.ErrorIs(t, c.TruncateCollections(ctx, testNilColl), ErrNilClient)
}

// TestTestClient_NilSafety tests method calls on nil TestClient instances.
func TestTestClient_NilSafety(t *testing.T) {
	var nilClient *TestClient

	verifyNilAccessors(t, nilClient)
	verifyNilLifecycle(t, nilClient)
	verifyNilFixtures(t, nilClient)
}

// TestTestClient_UninitializedSafety tests method calls on uninitialized TestClient.
func TestTestClient_UninitializedSafety(t *testing.T) {
	uninitializedClient := &TestClient{
		rawClient:       nil,
		defaultDatabase: testNilDB,
	}

	verifyNilAccessors(t, uninitializedClient)
	verifyNilLifecycle(t, uninitializedClient)
	verifyNilFixtures(t, uninitializedClient)
}

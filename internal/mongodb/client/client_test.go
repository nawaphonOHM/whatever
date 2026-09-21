package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Constants for test database and collection fixtures.
const (
	testMongoURI   = "mongodb://localhost:27017"
	testDefaultDB  = "default_db"
	testCustomDB   = "custom_db"
	testCollUsers  = "users"
	testCollOrders = "orders"
)

// createTestRawClient creates a connected mongo.Client and cleanup callback.
func createTestRawClient(t *testing.T) (*mongo.Client, func()) {
	// Connect to local test instance.
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	require.NotNil(t, rawClient)
	// Disconnect upon test completion.
	cleanup := func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}
	return rawClient, cleanup
}

// verifyDefaultAccessors checks database and collection resolutions.
func verifyDefaultAccessors(t *testing.T, client *Client) {
	// Verify default database accessor resolution.
	dbDefault := client.Database()
	require.NotNil(t, dbDefault)
	assert.Equal(t, testDefaultDB, dbDefault.Name())
	// Verify custom database override accessor resolution.
	dbCustom := client.Database(testCustomDB)
	require.NotNil(t, dbCustom)
	assert.Equal(t, testCustomDB, dbCustom.Name())
	// Verify collection resolution under default database.
	collDefault := client.Collection(testCollUsers)
	require.NotNil(t, collDefault)
	assert.Equal(t, testCollUsers, collDefault.Name())
	assert.Equal(t, testDefaultDB, collDefault.Database().Name())
}

// TestNewClient_Accessors tests database and collection accessors.
func TestNewClient_Accessors(t *testing.T) {
	// Initialize client instance for accessor tests.
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	client := NewClient(rawClient, testDefaultDB)
	require.NotNil(t, client)
	verifyDefaultAccessors(t, client)

	// Verify custom collection and raw client reference.
	collCustom := client.Collection(testCollOrders, testCustomDB)
	require.NotNil(t, collCustom)
	assert.Equal(t, testCollOrders, collCustom.Name())
	assert.Equal(t, testCustomDB, collCustom.Database().Name())
	assert.Same(t, rawClient, client.RawClient())
}

// verifyEmptyDefaultAccessors checks empty database accessor behavior.
func verifyEmptyDefaultAccessors(t *testing.T, client *Client) {
	// Default database should return an empty name.
	dbDefault := client.Database()
	require.NotNil(t, dbDefault)
	assert.Empty(t, dbDefault.Name())
	// Custom database should still resolve correctly.
	dbCustom := client.Database(testCustomDB)
	require.NotNil(t, dbCustom)
	assert.Equal(t, testCustomDB, dbCustom.Name())
	// Default collection database should have an empty name.
	collDefault := client.Collection(testCollUsers)
	require.NotNil(t, collDefault)
	assert.Equal(t, testCollUsers, collDefault.Name())
	assert.Empty(t, collDefault.Database().Name())
}

// TestNewClient_EmptyDefaultDatabase_Accessors tests empty default database.
func TestNewClient_EmptyDefaultDatabase_Accessors(t *testing.T) {
	// Initialize client with empty default database.
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	client := NewClient(rawClient, "")
	require.NotNil(t, client)
	verifyEmptyDefaultAccessors(t, client)
}

package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testMongoURI   = "mongodb://localhost:28018"
	testDefaultDB  = "default_db"
	testCustomDB   = "custom_db"
	testCollUsers  = "users"
	testCollOrders = "orders"
)

// createTestRawClient creates a test raw mongo.Client.
func createTestRawClient(t *testing.T) (*mongo.Client, func()) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	require.NotNil(t, rawClient)
	cleanup := func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}
	return rawClient, cleanup
}

// verifyAccessors checks database and collection resolutions.
func verifyAccessors(t *testing.T, tc *TestClient) {
	dbDefault := tc.Database()
	require.NotNil(t, dbDefault)
	assert.Equal(t, testDefaultDB, dbDefault.Name())

	dbCustom := tc.Database(testCustomDB)
	require.NotNil(t, dbCustom)
	assert.Equal(t, testCustomDB, dbCustom.Name())

	collDefault := tc.Collection(testCollUsers)
	require.NotNil(t, collDefault)
	assert.Equal(t, testCollUsers, collDefault.Name())
	assert.Equal(t, testDefaultDB, collDefault.Database().Name())
}

// verifyCustomAccessors checks custom collection and raw client handles.
func verifyCustomAccessors(t *testing.T, tc *TestClient, rawClient *mongo.Client) {
	collCustom := tc.Collection(testCollOrders, testCustomDB)
	require.NotNil(t, collCustom)
	assert.Equal(t, testCollOrders, collCustom.Name())
	assert.Equal(t, testCustomDB, collCustom.Database().Name())
	assert.Same(t, rawClient, tc.RawClient())
	assert.NotNil(t, tc.Client())
}

// TestNewClient_Accessors tests database and collection handle resolution.
func TestNewClient_Accessors(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	tc := NewClient(rawClient, testDefaultDB)
	require.NotNil(t, tc)
	verifyAccessors(t, tc)
	verifyCustomAccessors(t, tc, rawClient)
}

// TestNewClient_EmptyDefaultDatabase tests accessors when default db is empty.
func TestNewClient_EmptyDefaultDatabase(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	tc := NewClient(rawClient, "")
	require.NotNil(t, tc)

	dbDefault := tc.Database()
	require.NotNil(t, dbDefault)
	assert.Empty(t, dbDefault.Name())

	dbCustom := tc.Database(testCustomDB)
	require.NotNil(t, dbCustom)
	assert.Equal(t, testCustomDB, dbCustom.Name())
}

// TestTestClient_Ping_CanceledContext tests ping with a canceled context.
func TestTestClient_Ping_CanceledContext(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	tc := NewClient(rawClient, testDefaultDB)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.ErrorIs(t, tc.Ping(ctx), context.Canceled)
}

// TestTestClient_Disconnect_Close tests disconnect and close methods.
func TestTestClient_Disconnect_Close(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)

	tc := NewClient(rawClient, testDefaultDB)
	assert.NoError(t, tc.Close(context.Background()))

	// Subsequent disconnect or close returns error since client is disconnected.
	assert.Error(t, tc.Disconnect(context.Background()))
	assert.Error(t, tc.Close())
}

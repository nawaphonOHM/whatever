package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func verifyFirestoreFlags(t *testing.T, rawClient *mongo.Client) {
	assert.False(t, NewClient(rawClient, testDefaultDB).IsFirestore())
	assert.False(t, NewClient(rawClient, testDefaultDB, false).IsFirestore())
	assert.True(t, NewClient(rawClient, testDefaultDB, true).IsFirestore())

	var nilClient *Client
	assert.False(t, nilClient.IsFirestore())
}

// TestNewClient_FirestoreFlag tests Firestore mode configuration and IsFirestore accessor.
func TestNewClient_FirestoreFlag(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	verifyFirestoreFlags(t, rawClient)
}

// TestClient_Ping_Firestore_MockSuccess tests that Ping invokes Firestore ping when configured.
func TestClient_Ping_Firestore_MockSuccess(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	var calledDB string
	var calledCount int
	restore := SetMockPingFirestore(func(ctx context.Context, rc *mongo.Client, dbName string) error {
		assert.NotNil(t, ctx)
		assert.NotNil(t, rc)
		calledCount++
		calledDB = dbName
		return nil
	})
	defer restore()

	client := NewClient(rawClient, testDefaultDB, true)
	require.NoError(t, client.Ping(context.Background()))
	assert.Equal(t, 1, calledCount)
	assert.Equal(t, testDefaultDB, calledDB)
}

// TestClient_Ping_Firestore_MockFailure tests error propagation from Firestore ping.
func TestClient_Ping_Firestore_MockFailure(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	expectedErr := assert.AnError
	restore := SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		return expectedErr
	})
	defer restore()

	client := NewClient(rawClient, testDefaultDB, true)
	require.ErrorIs(t, client.Ping(context.Background()), expectedErr)
}

func setupFirestorePingTracker() (*int, func()) {
	var count int
	restore := SetMockPingFirestore(func(context.Context, *mongo.Client, string) error {
		count++
		return nil
	})
	return &count, restore
}

// TestClient_Ping_Routing verifies routing between standard MongoDB ping and Firestore ping.
func TestClient_Ping_Routing(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	called, restore := setupFirestorePingTracker()
	defer restore()

	assert.False(t, NewClient(rawClient, testDefaultDB, false).IsFirestore())
	clientFirestore := NewClient(rawClient, testDefaultDB, true)
	assert.True(t, clientFirestore.IsFirestore())
	require.NoError(t, clientFirestore.Ping(context.Background()))
	assert.Equal(t, 1, *called)
}

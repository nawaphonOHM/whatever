package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testUnreachableURI = "mongodb://127.0.0.1:59999"
	testShortDur       = 50 * time.Millisecond
	testContextDur     = 100 * time.Millisecond
)

// TestClient_Disconnect tests single and duplicate Disconnect calls.
func TestClient_Disconnect(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)

	client := NewClient(rawClient, testDefaultDB)
	require.NotNil(t, client)

	ctx := context.Background()
	assert.NoError(t, client.Disconnect(ctx))

	// Second disconnect returns an error because client is disconnected
	err = client.Disconnect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disconnected")
}

// TestClient_Ping_CanceledContext tests ping with a canceled context.
func TestClient_Ping_CanceledContext(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}()

	client := NewClient(rawClient, testDefaultDB)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.ErrorIs(t, client.Ping(ctx), context.Canceled)
}

// TestClient_Ping_Unreachable tests ping against an unreachable cluster.
func TestClient_Ping_Unreachable(t *testing.T) {
	rawClient, err := mongo.Connect(
		options.Client().
			ApplyURI(testUnreachableURI).
			SetServerSelectionTimeout(testShortDur).
			SetTimeout(testShortDur),
	)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}()

	client := NewClient(rawClient, testDefaultDB)
	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	assert.Error(t, client.Ping(ctx))
}

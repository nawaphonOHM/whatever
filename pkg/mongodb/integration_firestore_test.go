package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/pkg/mongodb"
)

const (
	envFirestoreURI       = "FIRESTORE_MONGODB_URI"
	envTestFirestoreURI   = "TEST_FIRESTORE_URI"
	testFirestoreDemoHost = "demo.asia-southeast1.firestore.goog"
	testFirestoreDemoDB   = "main"
	testLiveConnTimeout   = 10 * time.Second
)

func assertFirestoreDocQuery(ctx context.Context, t *testing.T, client *mongodb.Client, name string) {
	t.Helper()
	coll := client.Collection(name)
	var doc bson.M
	err := coll.FindOne(ctx, bson.D{}).Decode(&doc)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		require.NoError(t, err)
	}
}

func assertLiveFirestoreClient(ctx context.Context, t *testing.T, client *mongodb.Client) {
	t.Helper()
	assert.True(t, client.IsFirestore())
	require.NoError(t, client.Ping(ctx))
	db := client.Database()
	require.NotNil(t, db)
	names, err := db.ListCollectionNames(ctx, bson.D{})
	require.NoError(t, err)
	if len(names) > 0 {
		assertFirestoreDocQuery(ctx, t, client, names[0])
	}
}

func runLiveFirestoreConnect(ctx context.Context, t *testing.T) {
	t.Helper()
	client, err := mongodb.Connect(ctx)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer func() {
		require.NoError(t, client.Disconnect(ctx))
	}()

	assertLiveFirestoreClient(ctx, t, client)
}

func TestConnect_Firestore_LiveIntegration(t *testing.T) {
	uri := getLiveFirestoreURI()
	if uri == "" {
		t.Skip("skipping live Firestore integration test: FIRESTORE_MONGODB_URI not configured")
	}
	setupLiveFirestoreEnv(t, uri)
	ctx, cancel := context.WithTimeout(context.Background(), testLiveConnTimeout)
	defer cancel()

	runLiveFirestoreConnect(ctx, t)
}

func TestConnect_Firestore_MockIntegration(t *testing.T) {
	setupFirestoreDemoEnv(t)
	var pingCalled int
	cleanup := setupMockFirestorePingProbe(func() {
		pingCalled++
	})
	defer cleanup()

	client, err := mongodb.Connect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.True(t, client.IsFirestore())
	assert.Equal(t, 1, pingCalled)
}

func TestConnect_Firestore_ResetFallbackIntegration(t *testing.T) {
	setupDefaultEnv(t)
	var attempts int
	cleanup := setupResetFallbackMocks(&attempts)
	defer cleanup()

	client, err := mongodb.Connect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 2, attempts)
}

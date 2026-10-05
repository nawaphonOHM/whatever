//go:build testcontainers

package firestore_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
	mongoclient "github.com/nawaphonOHM/whatever/v2/pkg/testing/mongodb"
)

const (
	envLiveFirestoreURI = "FIRESTORE_MONGODB_URI"
	envTestFirestoreURI = "TEST_FIRESTORE_URI"
	testLiveTimeout     = 10 * time.Second
)

func getFirestoreLiveURI() string {
	if uri := os.Getenv(envLiveFirestoreURI); uri != "" {
		return uri
	}
	return os.Getenv(envTestFirestoreURI)
}

func assertLiveWireDocQuery(ctx context.Context, t *testing.T, client *mongoclient.TestClient, collName string) {
	t.Helper()
	coll := client.Collection(collName)
	var doc bson.M
	err := coll.FindOne(ctx, bson.D{}).Decode(&doc)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		require.NoError(t, err)
	}
}

func assertLiveWireQueries(ctx context.Context, t *testing.T, client *mongoclient.TestClient) {
	t.Helper()
	names, err := client.ListCollectionNames(ctx)
	require.NoError(t, err)
	if len(names) > 0 {
		assertLiveWireDocQuery(ctx, t, client, names[0])
	}
}

func runLiveWireClient(ctx context.Context, t *testing.T, uri string) {
	t.Helper()
	client, err := mongoclient.ConnectURI(ctx, uri)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer func() {
		require.NoError(t, client.Close())
	}()

	assert.True(t, client.IsFirestore())
	require.NoError(t, client.Ping(ctx))
	assertLiveWireQueries(ctx, t, client)
}

func TestFirestore_Integration_LiveMongoDB(t *testing.T) {
	uri := getFirestoreLiveURI()
	if uri == "" {
		t.Skip("skipping live Firestore MongoDB wire test: FIRESTORE_MONGODB_URI not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), testLiveTimeout)
	defer cancel()

	runLiveWireClient(ctx, t, uri)
}

func TestFirestore_Integration_ConfigAndProbeRules(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Host = "fc28fdab-9482-4359-9d76-0bb5f9fd1c4d.asia-southeast1.firestore.goog"
	cfg.Database = "main"

	assert.True(t, cfg.IsFirestore())
	builtURI := cfg.BuildURI(false)
	assert.Contains(t, builtURI, ":443")
	assert.Contains(t, builtURI, "tls=true")
	assert.Contains(t, builtURI, "loadBalanced=true")
	assert.Contains(t, builtURI, "retryWrites=false")
	assert.Contains(t, builtURI, "authMechanism=SCRAM-SHA-256")
}

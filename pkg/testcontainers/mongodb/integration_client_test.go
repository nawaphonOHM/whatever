//go:build testcontainers

package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	mongoclient "github.com/nawaphonOHM/whatever/v2/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/mongodb"
)

func assertClientPing(ctx context.Context, t *testing.T, c *mongodb.Container) *mongoclient.Client {
	t.Helper()
	client, err := c.Client(ctx)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NoError(t, client.Ping(ctx))
	return client
}

func assertInsertDoc(ctx context.Context, t *testing.T, client *mongoclient.Client, doc testDoc) {
	t.Helper()
	coll := client.Collection(testColName)
	require.NotNil(t, coll)
	res, err := coll.InsertOne(ctx, doc)
	require.NoError(t, err)
	assert.NotNil(t, res.InsertedID)
}

func assertQueryDoc(ctx context.Context, t *testing.T, client *mongoclient.Client, expected testDoc) {
	t.Helper()
	coll := client.Collection(testColName)
	var found testDoc
	filter := bson.D{{Key: "key", Value: expected.Key}}
	err := coll.FindOne(ctx, filter).Decode(&found)
	require.NoError(t, err)
	assert.Equal(t, expected.Key, found.Key)
	assert.Equal(t, expected.Value, found.Value)
}

func assertDropDatabase(ctx context.Context, t *testing.T, client *mongoclient.Client) {
	t.Helper()
	db := client.Database()
	require.NotNil(t, db)
	require.NoError(t, db.Drop(ctx))
}

func TestContainer_Integration_ClientOperations(t *testing.T) {
	requireSharedReady(t)
	ctx := context.Background()
	client := assertClientPing(ctx, t, sharedContainer)
	doc := testDoc{ID: "doc-1", Key: "env", Value: "production"}
	assertInsertDoc(ctx, t, client, doc)
	assertQueryDoc(ctx, t, client, doc)
	assertDropDatabase(ctx, t, client)
}

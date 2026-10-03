package mongodb

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testLifecycleDB   = "lifecycle_db"
	testLifecycleColl = "lifecycle_users"
	testConcurrentRun = 10
)

// TestSystemCollectionFilter tests detection and filtering of system collections.
func TestSystemCollectionFilter(t *testing.T) {
	assert.True(t, isSystemCollection("system.views"))
	assert.True(t, isSystemCollection("system.indexes"))
	assert.False(t, isSystemCollection("users"))
	assert.False(t, isSystemCollection("orders"))

	input := []string{"users", "system.views", "orders", "system.profile"}
	expected := []string{"users", "orders"}
	assert.Equal(t, expected, filterNonSystemCollections(input))
}

// verifyLifecycleCanceled checks fixture operations with a canceled context.
func verifyLifecycleCanceled(ctx context.Context, t *testing.T, tc *TestClient) {
	assert.Error(t, tc.DropDatabase(ctx))
	assert.Error(t, tc.DropDatabase(ctx, testLifecycleDB))
	assert.Error(t, tc.DropCollection(ctx, testLifecycleColl))
	assert.Error(t, tc.DropCollection(ctx, testLifecycleColl, testLifecycleDB))

	names, err := tc.ListCollectionNames(ctx)
	assert.Nil(t, names)
	assert.Error(t, err)

	assert.Error(t, tc.TruncateCollections(ctx, testLifecycleColl))
	assert.Error(t, tc.TruncateCollections(ctx))
}

// TestClientLifecycle_CanceledContext tests fixture operations against canceled contexts.
func TestClientLifecycle_CanceledContext(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}()

	tc := NewClient(rawClient, testLifecycleDB)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	verifyLifecycleCanceled(ctx, t, tc)
}

// TestClientLifecycle_TruncateEmptyName tests truncating empty named collection.
func TestClientLifecycle_TruncateEmptyName(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}()

	tc := NewClient(rawClient, testLifecycleDB)
	assert.NoError(t, tc.truncateNamed(context.Background(), ""))
}

// runConcurrentAccessors executes accessor methods in a goroutine.
func runConcurrentAccessors(t *testing.T, tc *TestClient) {
	assert.NotNil(t, tc.Database())
	assert.NotNil(t, tc.Collection(testLifecycleColl))
	assert.NotNil(t, tc.RawClient())
	assert.NotNil(t, tc.Client())
}

// TestTestClient_Concurrency tests thread safety of TestClient accessors.
func TestTestClient_Concurrency(t *testing.T) {
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, rawClient.Disconnect(context.Background()))
	}()

	tc := NewClient(rawClient, testLifecycleDB)
	var wg sync.WaitGroup
	for range testConcurrentRun {
		wg.Go(func() {
			runConcurrentAccessors(t, tc)
		})
	}
	wg.Wait()
}

package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
)

func TestChooseTargetColl_ReturnsProbe(t *testing.T) {
	assert.Equal(t, "probe", chooseTargetColl(nil))
	assert.Equal(t, "probe", chooseTargetColl(make([]string, 0)))
	assert.Equal(t, "users", chooseTargetColl([]string{"users"}))
}

func TestVerifyClientProbe_PropagatesFirestoreFlag_True(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	var capturedIsFirestore bool
	restoreProbe := SetMockProbe(func(ctx context.Context, c *client.Client, coll string) error {
		_ = ctx
		_ = coll
		capturedIsFirestore = c.IsFirestore()
		return nil
	})
	defer restoreProbe()

	err := verifyClientProbe(context.Background(), rawClient, testDBName, true)
	require.NoError(t, err)
	assert.True(t, capturedIsFirestore)
}

func TestVerifyClientProbe_PropagatesFirestoreFlag_False(t *testing.T) {
	rawClient, cleanup := createTestRawClient(t)
	defer cleanup()

	var capturedIsFirestore bool
	restoreProbe := SetMockProbe(func(ctx context.Context, c *client.Client, coll string) error {
		_ = ctx
		_ = coll
		capturedIsFirestore = c.IsFirestore()
		return nil
	})
	defer restoreProbe()

	err := verifyClientProbe(context.Background(), rawClient, testDBName, false)
	require.NoError(t, err)
	assert.False(t, capturedIsFirestore)
}

func TestDefaultPingFirestoreClient_Disconnected(t *testing.T) {
	rawClient, _ := getDisconnectedDatabase(t)
	err := defaultPingFirestoreClient(context.Background(), rawClient, testDBName)
	assert.Error(t, err)
}

func TestProbeFirestoreCollections_Disconnected(t *testing.T) {
	_, db := getDisconnectedDatabase(t)
	colls, err := probeFirestoreCollections(context.Background(), db)
	assert.Nil(t, colls)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list collections without maxTimeMS")
}

func TestDefaultProbeClient_FirestoreRouting(t *testing.T) {
	rawClient, _ := getDisconnectedDatabase(t)
	c := client.NewClient(rawClient, testDBName, true)
	err := defaultProbeClient(context.Background(), c, testDBName)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list collections without maxTimeMS")
}

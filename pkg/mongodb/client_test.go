package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/pkg/mongodb"
)

func TestConnect_InvalidConfig(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "-5s")

	ctx := context.Background()
	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connect timeout cannot be negative")
	assert.Nil(t, client)
}

func TestConnect_PingFailure(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://127.0.0.1:59999")
	t.Setenv("OHM9969_MONGODB_DATABASE", "testdb")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "50ms")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "50ms")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "50ms")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, mongodb.ErrNilClient)
	assert.Equal(t, "mongodb client is not initialized", mongodb.ErrNilClient.Error())

	assert.NotNil(t, mongodb.ErrNilConfig)
	assert.Equal(t, "mongodb config cannot be nil", mongodb.ErrNilConfig.Error())
}

func TestConnect_CanceledContext(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://localhost:27017")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

func TestNilClient(t *testing.T) {
	var client *mongodb.Client
	ctx := context.Background()

	assert.Nil(t, client.Database())
	assert.Nil(t, client.Database("test"))
	assert.Nil(t, client.Collection("users"))
	assert.Nil(t, client.Collection("users", "test"))
	assert.Nil(t, client.RawClient())
	assert.ErrorIs(t, client.Ping(ctx), mongodb.ErrNilClient)
	assert.ErrorIs(t, client.Disconnect(ctx), mongodb.ErrNilClient)
}

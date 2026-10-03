package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/pkg/testing/mongodb"
)

const (
	testInvalidPort     = 70000
	testUnreachablePort = 59999
	testTimeoutDuration = 100 * time.Millisecond
)

func TestConnect_InvalidPort(t *testing.T) {
	ctx := context.Background()
	client, err := mongodb.Connect(ctx, mongodb.WithPort(testInvalidPort))
	assert.Nil(t, client)
	require.Error(t, err)
	assert.ErrorIs(t, err, mongodb.ErrInvalidPort)
}

func TestConnect_PingDisabled(t *testing.T) {
	ctx := context.Background()
	client, err := mongodb.Connect(ctx, mongodb.WithPing(false), mongodb.WithDatabase("test_db"))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NotNil(t, client.RawClient())
	assert.NotNil(t, client.Database())
	assert.NotNil(t, client.Collection("users"))
	require.NoError(t, client.Close())
}

func TestConnect_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeoutDuration)
	cancel()
	client, err := mongodb.Connect(ctx, mongodb.WithPort(testUnreachablePort))
	assert.Nil(t, client)
	require.Error(t, err)
}

func TestConnectURI_Empty(t *testing.T) {
	ctx := context.Background()
	client, err := mongodb.ConnectURI(ctx, "")
	assert.Nil(t, client)
	require.ErrorIs(t, err, mongodb.ErrEmptyURI)

	client, err = mongodb.ConnectURI(ctx, "   ")
	assert.Nil(t, client)
	require.ErrorIs(t, err, mongodb.ErrEmptyURI)
}

func TestConnectURI_PingDisabled(t *testing.T) {
	ctx := context.Background()
	client, err := mongodb.ConnectURI(ctx, "mongodb://localhost:27017/custom_db", mongodb.WithPing(false))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.NotNil(t, client.Database())
	assert.Equal(t, "custom_db", client.Database().Name())
	require.NoError(t, client.Close())
}

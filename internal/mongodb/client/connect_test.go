package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConnUser = "testuser"
	testConnPass = "testpass"
	testConnDB   = "testdb"
)

// setupConnectPingEnv configures environment variables targeting an
// unreachable port.
func setupConnectPingEnv(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_URI", testUnreachableURI)
	t.Setenv("OHM9996_MONGODB_DATABASE", testConnDB)
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)
	t.Setenv("OHM9996_MONGODB_CONNECT_TIMEOUT", "50ms")
	t.Setenv("OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT", "50ms")
	t.Setenv("OHM9996_MONGODB_SOCKET_TIMEOUT", "50ms")
}

// TestConnect_LoadConfigFailure tests failure when config cannot be parsed.
func TestConnect_LoadConfigFailure(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_URI", testUnreachableURI)
	t.Setenv("OHM9996_MONGODB_DATABASE", testConnDB)
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)
	t.Setenv("OHM9996_MONGODB_CONNECT_TIMEOUT", "-5s")

	ctx := context.Background()
	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connect timeout cannot be negative")
	assert.Nil(t, client)
}

// TestConnect_PingFailure tests connection timeout with environment variables.
func TestConnect_PingFailure(t *testing.T) {
	setupConnectPingEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

// TestConnect_CanceledContext tests connection with pre-canceled context.
func TestConnect_CanceledContext(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_URI", testMongoURI)
	t.Setenv("OHM9996_MONGODB_DATABASE", testConnDB)
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

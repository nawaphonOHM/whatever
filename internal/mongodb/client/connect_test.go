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
	testPort59   = "59999"
)

// setupConnectPingEnv configures environment variables targeting an
// unreachable port.
func setupConnectPingEnv(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_HOST", "127.0.0.1")
	t.Setenv("OHM9996_MONGODB_PORT", testPort59)
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)
}

func setupConnectCanceledEnv(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_HOST", "localhost")
	t.Setenv("OHM9996_MONGODB_PORT", "28018")
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)
}

// TestConnect_LoadConfigFailure tests failure when config cannot be parsed.
func TestConnect_LoadConfigFailure(t *testing.T) {
	t.Setenv("OHM9996_MONGODB_HOST", "127.0.0.1")
	t.Setenv("OHM9996_MONGODB_PORT", "invalid-port")
	t.Setenv("OHM9996_MONGODB_USERNAME", testConnUser)
	t.Setenv("OHM9996_MONGODB_PASSWORD", testConnPass)

	ctx := context.Background()
	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load mongodb config")
	assert.Nil(t, client)
}

// TestConnect_PingFailure tests connection timeout with environment variables.
func TestConnect_PingFailure(t *testing.T) {
	setupConnectPingEnv(t)
	exitCalled, cleanup := setupExitCapture()
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

// TestConnect_CanceledContext tests connection with pre-canceled context.
func TestConnect_CanceledContext(t *testing.T) {
	setupConnectCanceledEnv(t)
	exitCalled, cleanup := setupExitCapture()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envMongoHost     = "OHM9996_MONGODB_HOST"
	envMongoPort     = "OHM9996_MONGODB_PORT"
	envMongoUsername = "OHM9996_MONGODB_USERNAME"
	envMongoPassword = "OHM9996_MONGODB_PASSWORD"
	envMongoDatabase = "OHM9996_MONGODB_DATABASE"
	testConnUser     = "testuser"
	testConnPass     = "testpass"
	testPort59       = "59999"
)

// setupConnectPingEnv configures environment variables targeting an
// unreachable port.
func setupConnectPingEnv(t *testing.T) {
	t.Setenv(envMongoHost, "127.0.0.1")
	t.Setenv(envMongoPort, testPort59)
	t.Setenv(envMongoUsername, testConnUser)
	t.Setenv(envMongoPassword, testConnPass)
}

func setupConnectCanceledEnv(t *testing.T) {
	t.Setenv(envMongoHost, "localhost")
	t.Setenv(envMongoPort, "28018")
	t.Setenv(envMongoUsername, testConnUser)
	t.Setenv(envMongoPassword, testConnPass)
}

// TestConnect_LoadConfigFailure tests failure when config cannot be parsed.
func TestConnect_LoadConfigFailure(t *testing.T) {
	t.Setenv(envMongoHost, "127.0.0.1")
	t.Setenv(envMongoPort, "invalid-port")
	t.Setenv(envMongoUsername, testConnUser)
	t.Setenv(envMongoPassword, testConnPass)

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

// TestConnect_ProbeFailure tests connection failure when probe fails after ping succeeds.
func TestConnect_ProbeFailure(t *testing.T) {
	setupConnectPingEnv(t)
	exitCalled, cleanupExit := setupExitCapture()
	defer cleanupExit()

	cleanup := setupMockPingAndProbe(assert.AnError)
	defer cleanup()

	client, err := Connect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to probe mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/pkg/mongodb"
)

// Constants for environment variable keys and test fixtures.
const (
	envHost          = "OHM9996_MONGODB_HOST"
	envPort          = "OHM9996_MONGODB_PORT"
	envUsername      = "OHM9996_MONGODB_USERNAME"
	envPassword      = "OHM9996_MONGODB_PASSWORD"
	testFailHost     = "127.0.0.1"
	testFailPort     = "59999"
	testDefaultHost  = "localhost"
	testDefaultPort  = "28018"
	testUser         = "testuser"
	testPassword     = "testpass"
	errPingSubstring = "failed to ping mongodb"
)

// setupPingFailureEnv configures environment variables for invalid target port.
func setupPingFailureEnv(t *testing.T) {
	t.Setenv(envHost, testFailHost)
	t.Setenv(envPort, testFailPort)
	t.Setenv(envUsername, testUser)
	t.Setenv(envPassword, testPassword)
}

func setupDefaultEnv(t *testing.T) {
	t.Setenv(envHost, testDefaultHost)
	t.Setenv(envPort, testDefaultPort)
	t.Setenv(envUsername, testUser)
	t.Setenv(envPassword, testPassword)
}

func setupExitCapture() (*bool, func()) {
	var exitCalled bool
	prev := mongodb.SetExitFunc(func(int) { exitCalled = true })
	return &exitCalled, func() { mongodb.SetExitFunc(prev) }
}

// TestConnect_InvalidConfig tests connection failure on invalid configuration.
func TestConnect_InvalidConfig(t *testing.T) {
	t.Setenv(envHost, testFailHost)
	t.Setenv(envPort, "invalid-port")
	t.Setenv(envUsername, testUser)
	t.Setenv(envPassword, testPassword)

	ctx := context.Background()
	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load mongodb config")
	assert.Nil(t, client)
}

// TestConnect_PingFailure tests connection timeout when MongoDB is unreachable.
func TestConnect_PingFailure(t *testing.T) {
	setupPingFailureEnv(t)
	exitCalled, cleanup := setupExitCapture()
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), errPingSubstring)
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

// TestConnect_CanceledContext verifies behavior when context is pre-canceled.
func TestConnect_CanceledContext(t *testing.T) {
	setupDefaultEnv(t)
	exitCalled, cleanup := setupExitCapture()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), errPingSubstring)
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

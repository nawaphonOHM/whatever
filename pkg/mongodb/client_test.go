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
	envConnectTimeout    = "OHM9969_MONGODB_CONNECT_TIMEOUT"
	envServerSelection   = "OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT"
	envSocketTimeout     = "OHM9969_MONGODB_SOCKET_TIMEOUT"
	envURI               = "OHM9969_MONGODB_URI"
	envDatabase          = "OHM9969_MONGODB_DATABASE"
	testFailURI          = "mongodb://127.0.0.1:59999"
	testDefaultURI       = "mongodb://localhost:27017"
	testDBName           = "testdb"
	testShortTimeout     = "50ms"
	errPingSubstring     = "failed to ping mongodb"
	expectedNilClientErr = "mongodb client is not initialized"
	expectedNilConfigErr = "mongodb config cannot be nil"
)

// setupPingFailureEnv configures environment variables for invalid target port.
func setupPingFailureEnv(t *testing.T) {
	t.Setenv(envURI, testFailURI)
	t.Setenv(envDatabase, testDBName)
	t.Setenv(envConnectTimeout, testShortTimeout)
	t.Setenv(envServerSelection, testShortTimeout)
	t.Setenv(envSocketTimeout, testShortTimeout)
}

// TestConnect_InvalidConfig tests connection failure on invalid configuration.
func TestConnect_InvalidConfig(t *testing.T) {
	t.Setenv(envConnectTimeout, "-5s")

	ctx := context.Background()
	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connect timeout cannot be negative")
	assert.Nil(t, client)
}

// TestConnect_PingFailure tests connection timeout when MongoDB is unreachable.
func TestConnect_PingFailure(t *testing.T) {
	setupPingFailureEnv(t)

	ctx, cancel := context.WithTimeout(
		context.Background(), 100*time.Millisecond,
	)
	defer cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), errPingSubstring)
	assert.Nil(t, client)
}

// TestSentinelErrors verifies exported sentinel error messages.
func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, mongodb.ErrNilClient)
	assert.Equal(
		t, expectedNilClientErr, mongodb.ErrNilClient.Error(),
	)

	assert.NotNil(t, mongodb.ErrNilConfig)
	assert.Equal(
		t, expectedNilConfigErr, mongodb.ErrNilConfig.Error(),
	)
}

// TestConnect_CanceledContext verifies behavior when context is pre-canceled.
func TestConnect_CanceledContext(t *testing.T) {
	t.Setenv(envURI, testDefaultURI)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client, err := mongodb.Connect(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), errPingSubstring)
	assert.Nil(t, client)
}

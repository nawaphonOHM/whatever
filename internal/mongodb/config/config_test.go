package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEnvHost                   = "OHM9996_MONGODB_HOST"
	testEnvPort                   = "OHM9996_MONGODB_PORT"
	testEnvProtocol               = "OHM9996_MONGODB_PROTOCOL"
	testEnvDatabase               = "OHM9996_MONGODB_DATABASE"
	testEnvUsername               = "OHM9996_MONGODB_USERNAME"
	testEnvPassword               = "OHM9996_MONGODB_PASSWORD"
	testEnvAuthSource             = "OHM9996_MONGODB_AUTH_SOURCE"
	testEnvAppName                = "OHM9996_MONGODB_APP_NAME"
	testEnvUUIDRepresentation     = "OHM9996_MONGODB_UUID_REPRESENTATION"
	testEnvConnectTimeout         = "OHM9996_MONGODB_CONNECT_TIMEOUT"
	testEnvServerSelectionTimeout = "OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT"
	testEnvSocketTimeout          = "OHM9996_MONGODB_SOCKET_TIMEOUT"
	testEnvMaxConnIdleTime        = "OHM9996_MONGODB_MAX_CONN_IDLE_TIME"
	testEnvMaxPoolSize            = "OHM9996_MONGODB_MAX_POOL_SIZE"
	testEnvMinPoolSize            = "OHM9996_MONGODB_MIN_POOL_SIZE"
	testCustomPort                = 27018
	testDefaultMaxPoolVal         = 100
	testDefaultMinPoolVal         = 5
	testCustomMaxPoolVal          = 200
	testCustomMinPoolVal          = 10
	testDefaultTimeoutSec         = 10
	testDefaultSelectSec          = 5
	testDefaultIdleMin            = 10
)

// setupEmptyEnv clears all mongodb environment variable overrides.
func setupEmptyEnv(t *testing.T) {
	for _, envVar := range []string{
		testEnvHost, testEnvPort, testEnvProtocol, testEnvDatabase,
		testEnvUsername, testEnvPassword, testEnvAuthSource, testEnvAppName,
		testEnvUUIDRepresentation, testEnvConnectTimeout,
		testEnvServerSelectionTimeout,
		testEnvSocketTimeout, testEnvMaxConnIdleTime, testEnvMaxPoolSize, testEnvMinPoolSize,
	} {
		t.Setenv(envVar, "")
	}
}

// verifyDefaultBaseFields asserts default connection and identity fields.
func verifyDefaultBaseFields(t *testing.T, cfg *Config) {
	assert.Empty(t, cfg.Host)
	assert.Equal(t, defaultPort, cfg.Port)
	assert.Equal(t, ProtocolMongoDB, cfg.Protocol)
	assert.Empty(t, cfg.Database)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.Empty(t, cfg.AuthSource)
	assert.Empty(t, cfg.AppName)
	assert.Equal(t, UUIDRepresentationUnspecified, cfg.UUIDRepresentation)
}

// verifyDefaultRuntimeFields asserts default timeout and pooling fields.
func verifyDefaultRuntimeFields(t *testing.T, cfg *Config) {
	assert.Equal(t, testDefaultTimeoutSec*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, testDefaultSelectSec*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, testDefaultTimeoutSec*time.Second, cfg.SocketTimeout)
	assert.Equal(t, testDefaultIdleMin*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, uint64(testDefaultMaxPoolVal), cfg.MaxPoolSize)
	assert.Equal(t, uint64(testDefaultMinPoolVal), cfg.MinPoolSize)
}

// TestDefaultConfig verifies default values created by DefaultConfig.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)
	verifyDefaultBaseFields(t, cfg)
	verifyDefaultRuntimeFields(t, cfg)
	assert.Error(t, cfg.Validate())
}

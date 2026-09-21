package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testCustomURI  = "mongodb://custom-host:27018"
	testAppSvc     = "order-service"
	testCustomUser = "testuser"
	testCustomPass = "testpass"
	testMax50      = uint64(50)
	testMin10      = uint64(10)
	testConn15     = 15 * time.Second
	testSelect8    = 8 * time.Second
	testSock25     = 25 * time.Second
	testIdle5      = 5 * time.Minute
)

// buildCustomTestConfig constructs a Config populated with custom parameters.
func buildCustomTestConfig() *config.Config {
	return &config.Config{
		URI:                    testCustomURI,
		Database:               testCollOrders,
		Username:               testCustomUser,
		Password:               testCustomPass,
		AppName:                testAppSvc,
		ConnectTimeout:         testConn15,
		ServerSelectionTimeout: testSelect8,
		SocketTimeout:          testSock25,
		MaxConnIdleTime:        testIdle5,
		MaxPoolSize:            testMax50,
		MinPoolSize:            testMin10,
	}
}

// verifyCustomDriverTimeouts asserts custom timeout options on client options.
func verifyCustomDriverTimeouts(
	t *testing.T,
	clientOpts *options.ClientOptions,
) {
	// Verify custom connect and server selection durations.
	assert.Equal(t, testConn15, *clientOpts.ConnectTimeout)
	assert.Equal(t, testSelect8, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, testSock25, *clientOpts.Timeout)
	assert.Equal(t, testIdle5, *clientOpts.MaxConnIdleTime)
}

// verifyCustomDriverAuth asserts custom credentials on client options.
func verifyCustomDriverAuth(t *testing.T, clientOpts *options.ClientOptions) {
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, testCustomUser, clientOpts.Auth.Username)
	assert.Equal(t, testCustomPass, clientOpts.Auth.Password)
}

// TestBuildClientOptions_CustomConfig tests options with explicit values.
func TestBuildClientOptions_CustomConfig(t *testing.T) {
	// Build client options from fully populated config.
	cfg := buildCustomTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	// Verify custom URI and metadata.
	assert.Equal(t, testCustomURI, clientOpts.GetURI())
	assert.Equal(t, testAppSvc, *clientOpts.AppName)

	// Verify credentials, pool, and timeouts.
	verifyCustomDriverAuth(t, clientOpts)
	assert.Equal(t, testMax50, *clientOpts.MaxPoolSize)
	assert.Equal(t, testMin10, *clientOpts.MinPoolSize)
	verifyCustomDriverTimeouts(t, clientOpts)
}

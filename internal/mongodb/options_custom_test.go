package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testCustomURI = "mongodb://custom-host:27018"
	testAppSvc    = "order-service"
	testMax50     = uint64(50)
	testMin10     = uint64(10)
	testConn15    = 15 * time.Second
	testSelect8   = 8 * time.Second
	testSock25    = 25 * time.Second
	testIdle5     = 5 * time.Minute
)

// buildCustomTestConfig constructs a Config populated with custom parameters.
func buildCustomTestConfig() *Config {
	return &Config{
		BaseFields: BaseFields{
			URI:      testCustomURI,
			Database: testCollOrders,
			AppName:  testAppSvc,
		},
		TimeoutFields: TimeoutFields{
			ConnectTimeout:         testConn15,
			ServerSelectionTimeout: testSelect8,
			SocketTimeout:          testSock25,
			MaxConnIdleTime:        testIdle5,
		},
		PoolFields: PoolFields{
			MaxPoolSize: testMax50,
			MinPoolSize: testMin10,
		},
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

// TestBuildClientOptions_CustomConfig tests options with explicit values.
func TestBuildClientOptions_CustomConfig(t *testing.T) {
	// Build client options from fully populated config.
	cfg := buildCustomTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	// Verify custom URI and metadata.
	assert.Equal(t, testCustomURI, clientOpts.GetURI())
	assert.Equal(t, testAppSvc, *clientOpts.AppName)

	// Verify pool and timeout limits.
	assert.Equal(t, testMax50, *clientOpts.MaxPoolSize)
	assert.Equal(t, testMin10, *clientOpts.MinPoolSize)
	verifyCustomDriverTimeouts(t, clientOpts)
}

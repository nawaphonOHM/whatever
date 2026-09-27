package client

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

func buildPopulatedTestConfig() *config.Config {
	return &config.Config{
		Host:                   "custom-host",
		Port:                   testOptionsPort,
		Username:               testOptionsUser,
		Password:               testOptionsPass,
		AuthSource:             testOptionsUser,
		AppName:                testOptionsAppName,
		Protocol:               config.ProtocolMongoDB,
		UUIDRepresentation:     config.UUIDRepresentationStandard,
		ConnectTimeout:         testConnSec * time.Second,
		ServerSelectionTimeout: testSelectSec * time.Second,
		SocketTimeout:          testSocketSec * time.Second,
		MaxConnIdleTime:        testIdleMin * time.Minute,
		MaxPoolSize:            testCustomMaxPool,
		MinPoolSize:            testCustomMinPool,
	}
}

func verifyCustomClientOptions(t *testing.T, clientOpts *options.ClientOptions) {
	assert.Equal(t, testOptionsAppName, *clientOpts.AppName)
	assert.Equal(t, testConnSec*time.Second, *clientOpts.ConnectTimeout)
	assert.Equal(t, testSelectSec*time.Second, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, testSocketSec*time.Second, *clientOpts.Timeout)
	assert.Equal(t, testIdleMin*time.Minute, *clientOpts.MaxConnIdleTime)
	assert.Equal(t, uint64(testCustomMaxPool), *clientOpts.MaxPoolSize)
	assert.Equal(t, uint64(testCustomMinPool), *clientOpts.MinPoolSize)
}

// TestBuildClientOptions_PopulatedConfig tests options with populated config.
func TestBuildClientOptions_PopulatedConfig(t *testing.T) {
	cfg := buildPopulatedTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	expectedURI := "mongodb://admin:secretpassword@custom-host:27018/" +
		"?uuidRepresentation=standard&tls=false&authSource=admin&appName=my-service"
	assert.Equal(t, expectedURI, clientOpts.GetURI())
	verifyCustomClientOptions(t, clientOpts)
}

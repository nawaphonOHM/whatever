package client

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testOptionsPort    = 27018
	testOptionsUser    = "admin"
	testOptionsPass    = "secretpassword"
	testOptionsAppName = "my-service"
	testDefaultMaxPool = 100
	testDefaultMinPool = 5
	testCustomMaxPool  = 200
	testCustomMinPool  = 10
	testConnSec        = 15
	testSelectSec      = 3
	testSocketSec      = 20
	testIdleMin        = 5
)

func verifyDefaultClientTimeouts(t *testing.T, clientOpts *options.ClientOptions) {
	assert.Equal(t, 10*time.Second, *clientOpts.ConnectTimeout)
	assert.Equal(t, 5*time.Second, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, *clientOpts.Timeout)
	assert.Equal(t, 10*time.Minute, *clientOpts.MaxConnIdleTime)
}

func verifyDefaultClientPool(t *testing.T, clientOpts *options.ClientOptions) {
	assert.Equal(t, uint64(testDefaultMaxPool), *clientOpts.MaxPoolSize)
	assert.Equal(t, uint64(testDefaultMinPool), *clientOpts.MinPoolSize)
}

// TestBuildClientOptions_DefaultConfig tests options built from DefaultConfig.
func TestBuildClientOptions_DefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	verifyDefaultClientTimeouts(t, clientOpts)
	verifyDefaultClientPool(t, clientOpts)
}

// TestBuildClientOptions_NilConfig tests options built when config is nil.
func TestBuildClientOptions_NilConfig(t *testing.T) {
	clientOpts := BuildClientOptions(nil)
	require.NotNil(t, clientOpts)
	verifyDefaultClientTimeouts(t, clientOpts)
	verifyDefaultClientPool(t, clientOpts)
}

func buildAuthHostTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = "custom-host"
	cfg.Port = testOptionsPort
	cfg.Username = testOptionsUser
	cfg.Password = testOptionsPass
	return cfg
}

func verifyAuthAndHosts(t *testing.T, clientOpts *options.ClientOptions) {
	assert.Contains(t, clientOpts.Hosts, "custom-host:27018")
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, testOptionsUser, clientOpts.Auth.Username)
	assert.Equal(t, testOptionsPass, clientOpts.Auth.Password)
}

// TestBuildClientOptions_AuthAndHosts verifies hosts and credentials on options.
func TestBuildClientOptions_AuthAndHosts(t *testing.T) {
	cfg := buildAuthHostTestConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	verifyAuthAndHosts(t, clientOpts)
}

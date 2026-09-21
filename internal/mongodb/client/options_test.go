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
	testDefaultMaxPool = uint64(100)
	testDefaultMax5    = uint64(5)
)

// verifyDefaultDriverPool verifies pool size options on client options.
func verifyDefaultDriverPool(t *testing.T, clientOpts *options.ClientOptions) {
	// Verify pool size boundaries.
	assert.Equal(t, testDefaultMaxPool, *clientOpts.MaxPoolSize)
	assert.Equal(t, testDefaultMax5, *clientOpts.MinPoolSize)
	assert.Equal(t, 10*time.Minute, *clientOpts.MaxConnIdleTime)
	assert.Nil(t, clientOpts.AppName)
}

// TestBuildClientOptions_DefaultConfig tests options built from DefaultConfig.
func TestBuildClientOptions_DefaultConfig(t *testing.T) {
	// Build client options from default config.
	cfg := config.DefaultConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	// Verify timeouts on driver client options.
	assert.Empty(t, clientOpts.GetURI())
	assert.Equal(t, 10*time.Second, *clientOpts.ConnectTimeout)
	assert.Equal(t, 5*time.Second, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, *clientOpts.Timeout)
	verifyDefaultDriverPool(t, clientOpts)
}

// TestBuildClientOptions_NilConfig tests options built when config is nil.
func TestBuildClientOptions_NilConfig(t *testing.T) {
	// Default options should be used when nil is supplied.
	clientOpts := BuildClientOptions(nil)
	require.NotNil(t, clientOpts)

	assert.Empty(t, clientOpts.GetURI())
	assert.Equal(t, 10*time.Second, *clientOpts.ConnectTimeout)
}

// verifyZeroDriverOptions asserts all optional driver options are nil.
func verifyZeroDriverOptions(t *testing.T, clientOpts *options.ClientOptions) {
	// Verify all optional parameters are unset.
	assert.Nil(t, clientOpts.ConnectTimeout)
	assert.Nil(t, clientOpts.ServerSelectionTimeout)
	assert.Nil(t, clientOpts.Timeout)
	assert.Nil(t, clientOpts.MaxPoolSize)
	assert.Nil(t, clientOpts.MinPoolSize)
	assert.Nil(t, clientOpts.MaxConnIdleTime)
	assert.Nil(t, clientOpts.AppName)
}

// TestBuildClientOptions_ZeroValues tests options with zero-valued config.
func TestBuildClientOptions_ZeroValues(t *testing.T) {
	// Only URI is set on zero-valued config.
	cfg := &config.Config{
		URI: "mongodb://custom-host:27017",
	}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	assert.Equal(t, "mongodb://custom-host:27017", clientOpts.GetURI())
	verifyZeroDriverOptions(t, clientOpts)
}

// TestBuildClientOptions_Credentials tests credential propagation on driver client options.
func TestBuildClientOptions_Credentials(t *testing.T) {
	cfg := &config.Config{
		URI:      "mongodb://auth-host:27017",
		Database: "auth_db",
		Username: "admin",
		Password: "secretpassword",
	}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)
	require.NotNil(t, clientOpts.Auth)
	assert.Equal(t, "admin", clientOpts.Auth.Username)
	assert.Equal(t, "secretpassword", clientOpts.Auth.Password)
}

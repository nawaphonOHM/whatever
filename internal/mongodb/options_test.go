package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testDefaultMax5 = uint64(5)
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
	cfg := DefaultConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	// Verify timeouts on driver client options.
	assert.Equal(t, DefaultURI, clientOpts.GetURI())
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

	assert.Equal(t, DefaultURI, clientOpts.GetURI())
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
	cfg := &Config{URI: DefaultURI}
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	assert.Equal(t, DefaultURI, clientOpts.GetURI())
	verifyZeroDriverOptions(t, clientOpts)
}

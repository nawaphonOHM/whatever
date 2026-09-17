package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestBuildClientOptions_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	assert.Equal(t, "mongodb://localhost:27017", clientOpts.GetURI())
	assert.Equal(t, 10*time.Second, *clientOpts.ConnectTimeout)
	assert.Equal(t, 5*time.Second, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, *clientOpts.Timeout)
	assert.Equal(t, uint64(100), *clientOpts.MaxPoolSize)
	assert.Equal(t, uint64(5), *clientOpts.MinPoolSize)
	assert.Equal(t, 10*time.Minute, *clientOpts.MaxConnIdleTime)
	assert.Nil(t, clientOpts.AppName)
}

func TestBuildClientOptions_NilConfig(t *testing.T) {
	clientOpts := BuildClientOptions(nil)
	require.NotNil(t, clientOpts)

	assert.Equal(t, "mongodb://localhost:27017", clientOpts.GetURI())
	assert.Equal(t, 10*time.Second, *clientOpts.ConnectTimeout)
}

func TestBuildClientOptions_CustomConfig(t *testing.T) {
	cfg := &Config{
		URI:                    "mongodb://custom-host:27018",
		Database:               "orders",
		ConnectTimeout:         15 * time.Second,
		ServerSelectionTimeout: 8 * time.Second,
		SocketTimeout:          25 * time.Second,
		MaxPoolSize:            50,
		MinPoolSize:            10,
		MaxConnIdleTime:        5 * time.Minute,
		AppName:                "order-service",
	}

	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	assert.Equal(t, "mongodb://custom-host:27018", clientOpts.GetURI())
	assert.Equal(t, 15*time.Second, *clientOpts.ConnectTimeout)
	assert.Equal(t, 8*time.Second, *clientOpts.ServerSelectionTimeout)
	assert.Equal(t, 25*time.Second, *clientOpts.Timeout)
	assert.Equal(t, uint64(50), *clientOpts.MaxPoolSize)
	assert.Equal(t, uint64(10), *clientOpts.MinPoolSize)
	assert.Equal(t, 5*time.Minute, *clientOpts.MaxConnIdleTime)
	assert.Equal(t, "order-service", *clientOpts.AppName)
}

func TestBuildClientOptions_WithExtraDriverOptions(t *testing.T) {
	cfg := DefaultConfig()
	extra1 := options.Client().SetDirect(true)
	extra2 := options.Client().SetAppName("override-app")

	clientOpts := BuildClientOptions(cfg, extra1, extra2)
	require.NotNil(t, clientOpts)

	assert.True(t, *clientOpts.Direct)
	assert.Equal(t, "override-app", *clientOpts.AppName)
}

func TestNewOptions_And_WithDriverOptions(t *testing.T) {
	driverOpt1 := options.Client().SetDirect(true)
	driverOpt2 := options.Client().SetAppName("test-app")

	opts := NewOptions(WithDriverOptions(driverOpt1, nil, driverOpt2), nil)
	require.NotNil(t, opts)
	require.Len(t, opts.DriverOptions, 2)
	assert.Same(t, driverOpt1, opts.DriverOptions[0])
	assert.Same(t, driverOpt2, opts.DriverOptions[1])
}

func TestBuildClientOptions_ZeroValues(t *testing.T) {
	cfg := &Config{
		URI: "mongodb://localhost:27017",
	}

	clientOpts := BuildClientOptions(cfg)
	require.NotNil(t, clientOpts)

	assert.Equal(t, "mongodb://localhost:27017", clientOpts.GetURI())
	assert.Nil(t, clientOpts.ConnectTimeout)
	assert.Nil(t, clientOpts.ServerSelectionTimeout)
	assert.Nil(t, clientOpts.Timeout)
	assert.Nil(t, clientOpts.MaxPoolSize)
	assert.Nil(t, clientOpts.MinPoolSize)
	assert.Nil(t, clientOpts.MaxConnIdleTime)
	assert.Nil(t, clientOpts.AppName)
}

func TestNewOptions_Empty(t *testing.T) {
	opts := NewOptions()
	require.NotNil(t, opts)
	assert.Empty(t, opts.DriverOptions)
}

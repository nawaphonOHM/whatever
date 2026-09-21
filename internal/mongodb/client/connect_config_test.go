package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testUnreachablePoolMax = uint64(10)
	testUnreachablePoolMin = uint64(1)
	testUnreachUser        = "unreach_user"
	testUnreachPass        = "unreach_pass"
)

// createUnreachableConfig returns a Config pointing to an unreachable port.
func createUnreachableConfig() *config.Config {
	return &config.Config{
		URI:                    testUnreachableURI,
		Database:               testDefaultDB,
		Username:               testUnreachUser,
		Password:               testUnreachPass,
		ConnectTimeout:         testShortDur,
		ServerSelectionTimeout: testShortDur,
		SocketTimeout:          testShortDur,
		MaxPoolSize:            testUnreachablePoolMax,
		MinPoolSize:            testUnreachablePoolMin,
	}
}

// TestConnectWithConfig_NilConfig tests rejection of nil config.
func TestConnectWithConfig_NilConfig(t *testing.T) {
	ctx := context.Background()
	client, err := ConnectWithConfig(ctx, nil)
	assert.ErrorIs(t, err, config.ErrNilConfig)
	assert.Nil(t, client)
}

// TestConnectWithConfig_InvalidConfig tests rejection of invalid config.
func TestConnectWithConfig_InvalidConfig(t *testing.T) {
	ctx := context.Background()
	invalidCfg := &config.Config{
		URI: "",
	}
	client, err := ConnectWithConfig(ctx, invalidCfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid mongodb config")
	assert.Nil(t, client)
}

// TestConnectWithConfig_InvalidURIFormat tests client creation failure.
func TestConnectWithConfig_InvalidURIFormat(t *testing.T) {
	ctx := context.Background()
	invalidURICfg := &config.Config{
		URI:            "://invalid uri",
		Database:       testDefaultDB,
		Username:       testUnreachUser,
		Password:       testUnreachPass,
		ConnectTimeout: 1 * time.Second,
	}
	client, err := ConnectWithConfig(ctx, invalidURICfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create mongodb client")
	assert.Nil(t, client)
}

// TestConnectWithConfig_PingFailure tests ping timeout against unreachable
// host.
func TestConnectWithConfig_PingFailure(t *testing.T) {
	cfg := createUnreachableConfig()
	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := ConnectWithConfig(ctx, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

// TestConnectWithConfig_WithOptions tests merging custom driver options.
func TestConnectWithConfig_WithOptions(t *testing.T) {
	cfg := createUnreachableConfig()
	extraOpt := options.Client().SetAppName("custom-connect-app")

	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := ConnectWithConfig(ctx, cfg, WithDriverOptions(extraOpt))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
}

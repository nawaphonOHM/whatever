package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

const (
	testUnreachUser = "unreach_user"
	testUnreachPass = "unreach_pass"
	testUnreachPort = 59999
)

// createUnreachableConfig returns a Config pointing to an unreachable port.
func createUnreachableConfig() *config.Config {
	return &config.Config{
		Host:               "127.0.0.1",
		Port:               testUnreachPort,
		Username:           testUnreachUser,
		Password:           testUnreachPass,
		Protocol:           config.ProtocolMongoDB,
		UUIDRepresentation: config.UUIDRepresentationUnspecified,
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
		Host: "",
	}
	client, err := ConnectWithConfig(ctx, invalidCfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid mongodb config")
	assert.Nil(t, client)
}

func interceptExitHook() (*bool, func()) {
	var exitCalled bool
	restore := SetExitFunc(func(int) {
		exitCalled = true
	})
	return &exitCalled, func() { SetExitFunc(restore) }
}

// TestConnectWithConfig_PingFailure tests ping timeout against unreachable host.
func TestConnectWithConfig_PingFailure(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := ConnectWithConfig(ctx, createUnreachableConfig())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

// TestConnectWithConfig_WithOptions tests merging custom driver options.
func TestConnectWithConfig_WithOptions(t *testing.T) {
	exitCalled, restore := interceptExitHook()
	defer restore()

	extraOpt := options.Client().SetAppName("custom-connect-app")
	ctx, cancel := context.WithTimeout(context.Background(), testContextDur)
	defer cancel()

	client, err := ConnectWithConfig(ctx, createUnreachableConfig(), WithDriverOptions(extraOpt))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to ping mongodb")
	assert.Nil(t, client)
	assert.True(t, *exitCalled)
}

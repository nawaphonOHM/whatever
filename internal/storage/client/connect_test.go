package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func newConnectTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = "connect-bucket"
	cfg.Endpoint = "https://s3.amazonaws.com"
	cfg.AccessKeyID = testAccessKey
	cfg.SecretAccessKey = testSecretKey
	return cfg
}

func assertConnectedClient(t *testing.T, c *client.Client, probed bool, bucket string) {
	t.Helper()
	assert.True(t, probed)
	assert.Equal(t, bucket, c.DefaultBucket())
	assert.Equal(t, contracts.ProviderAWS, c.Provider())
}

func TestConnectWithConfig_SuccessAndProbe(t *testing.T) {
	t.Parallel()

	probed := false
	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodHead {
			probed = true
		}
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})

	cfg := newConnectTestConfig()
	c, err := client.ConnectWithConfig(context.Background(), cfg, client.WithHTTPClient(httpClient))
	require.NoError(t, err)
	require.NotNil(t, c)
	assertConnectedClient(t, c, probed, "connect-bucket")
}

func TestConnectWithConfig_SkipProbe(t *testing.T) {
	t.Parallel()

	probed := false
	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		probed = true
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})

	cfg := newConnectTestConfig()
	c, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.False(t, probed)
}

func TestConnectWithConfig_Errors(t *testing.T) {
	t.Parallel()

	t.Run("nil config returns ErrNilConfig", func(t *testing.T) {
		t.Parallel()
		c, err := client.ConnectWithConfig(context.Background(), nil)
		assert.Nil(t, c)
		assert.ErrorIs(t, err, contracts.ErrNilConfig)
	})

	t.Run("invalid config returns error", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{Bucket: ""}
		c, err := client.ConnectWithConfig(context.Background(), cfg)
		assert.Nil(t, c)
		assert.Error(t, err)
	})
}

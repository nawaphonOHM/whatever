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

type capturedRequest struct {
	Host   string
	Path   string
	Scheme string
	Auth   string
}

func captureRequest(dst *capturedRequest) func(*http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		dst.Host = req.URL.Host
		dst.Path = req.URL.Path
		dst.Scheme = req.URL.Scheme
		dst.Auth = req.Header.Get("Authorization")
		return newHTTPResponse(http.StatusOK, "", nil), nil
	}
}

func testConfigWithKeys() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.AccessKeyID = testAccessKey
	cfg.SecretAccessKey = testSecretKey
	return cfg
}

func headWithConfig(t *testing.T, cfg *config.Config, got *capturedRequest) {
	t.Helper()
	httpClient := newMockHTTPClient(captureRequest(got))
	c, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.NoError(t, err)
	_, err = c.HeadObject(context.Background(), contracts.HeadObjectInput{Key: testObjectKey})
	require.NoError(t, err)
}

func assertSignedRequest(t *testing.T, got capturedRequest, host, path string) {
	t.Helper()
	assert.Equal(t, host, got.Host)
	assert.Equal(t, path, got.Path)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

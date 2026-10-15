package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
)

func TestConnect_Routing_BucketURL_ProxyPrefix(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://minio.internal.net/s3/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "minio.internal.net", got.Host)
	assert.Equal(t, "/s3/mybucket/k", got.Path)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

func TestConnect_Routing_BucketURL_TrailingPrefix(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://storage.googleapis.com/mybucket/some/prefix"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "storage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

func TestConnect_Routing_BucketURL_HostMismatch(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "mybucket"
	cfg.BucketURL = "https://other.storage.googleapis.com"
	httpClient := newMockHTTPClient(func(*http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})
	_, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

func TestConnect_Routing_BucketURL_PathMismatch(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "mybucket"
	cfg.BucketURL = "https://storage.googleapis.com/other"
	httpClient := newMockHTTPClient(func(*http.Request) (*http.Response, error) {
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})
	_, err := client.ConnectWithConfig(
		context.Background(),
		cfg,
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
}

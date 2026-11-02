package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
)

func TestConnect_Routing_AWSURLDerivedRegion(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.ap-southeast-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.ap-southeast-1.amazonaws.com", testVirtualPath)
	assert.Contains(t, got.Auth, "ap-southeast-1")
}

func TestConnect_Routing_AWSURLDerivedRegionTrailingDot(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.eu-west-1.amazonaws.com."
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.eu-west-1.amazonaws.com", testVirtualPath)
	assert.Contains(t, got.Auth, "eu-west-1")
}

func TestConnect_Routing_AWSBucketURLDerivedRegion(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.ap-southeast-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.ap-southeast-1.amazonaws.com", testVirtualPath)
	assert.Contains(t, got.Auth, "ap-southeast-1")
}

func TestConnect_Routing_AWSBucketURLDerivedRegionTrailingDot(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.eu-west-1.amazonaws.com."
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.eu-west-1.amazonaws.com", testVirtualPath)
	assert.Contains(t, got.Auth, "eu-west-1")
}

func TestConnect_Routing_AWSRegionConflict(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Region = "us-west-2"
	cfg.Endpoint = "https://s3.ap-southeast-1.amazonaws.com"
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
	assert.Contains(t, err.Error(), "region mismatch")
}

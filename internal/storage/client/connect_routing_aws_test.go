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

func TestConnect_Routing_AWSBucketURLDerivedRegion(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.ap-southeast-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.ap-southeast-1.amazonaws.com", testVirtualPath)
	assert.Contains(t, got.Auth, "ap-southeast-1")
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

func TestConnect_Routing_AWSPrivateLinkEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

func TestConnect_Routing_AWSFIPSEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3-fips.us-east-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3-fips.us-east-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

func TestConnect_Routing_AWSThreePartEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://mybucket.s3.us-east-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.us-east-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

func TestConnect_Routing_AWSInsecureEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.amazonaws.com"
	cfg.Insecure = true
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "http", got.Scheme)
	assert.Contains(t, got.Host, "s3.amazonaws.com")
}

func TestConnect_Routing_AWSCustomPortEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.amazonaws.com:8443"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Contains(t, got.Host, ":8443")
}

func TestConnect_Routing_AWSPathPrefixEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.amazonaws.com/proxy/s3"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Contains(t, got.Path, "/proxy/s3")
}

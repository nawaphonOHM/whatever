package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

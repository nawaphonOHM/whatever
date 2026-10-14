package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_Routing_AmbiguousHostAndPathBucket(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "data"
	cfg.BucketURL = "https://data.company.com/data"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "data.company.com", got.Host)
	assert.Equal(t, "/data/k", got.Path)
}

func TestConnect_Routing_GCSEndpointWithBucketHost(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "mybucket"
	cfg.Endpoint = "https://mybucket.storage.googleapis.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "storage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

func TestConnect_Routing_AWSWebsiteBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3-website-us-east-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.us-east-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "us-east-1")
}

func TestConnect_Routing_CustomVHostMinioBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "mybucket"
	cfg.BucketURL = "https://mybucket.minio.internal.net"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "minio.internal.net", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

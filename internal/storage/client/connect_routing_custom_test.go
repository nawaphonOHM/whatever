package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_Routing_VHostBucketURLWithPrefixAWS(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.eu-west-1.amazonaws.com/folder/sub"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.eu-west-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "eu-west-1")
}

func TestConnect_Routing_VHostBucketURLWithPrefixGCS(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.storage.googleapis.com/folder/"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "storage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

func TestConnect_Routing_CustomEndpointHostMatchesBucket(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "minio"
	cfg.Endpoint = "https://minio.default.svc.cluster.local:9000"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "minio.default.svc.cluster.local:9000", got.Host)
	assert.Equal(t, "/minio/k", got.Path)
}

func TestConnect_Routing_CustomEndpointProxyPathMatchesBucket(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Bucket = "data"
	cfg.Endpoint = "https://proxy.company.com/data/s3"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "proxy.company.com", got.Host)
	assert.Equal(t, "/data/s3/data/k", got.Path)
}

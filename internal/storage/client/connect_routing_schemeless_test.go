package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_Routing_SchemelessAWSEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "s3.ap-southeast-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.ap-southeast-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "ap-southeast-1")
}

func TestConnect_Routing_SchemelessAWSBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "mybucket.s3.eu-west-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.eu-west-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "eu-west-1")
}

func TestConnect_Routing_SchemelessGCSPathBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "storage.googleapis.com/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "storage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

func TestConnect_Routing_SchemelessGCSHostBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "mybucket.storage.googleapis.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "storage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

func TestConnect_Routing_SchemelessMinioBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "minio.internal.net/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "minio.internal.net", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_Routing_AWSRegional(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Region = testAWSRegionSE
	cfg.Endpoint = "https://s3.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, testAWSRegionalHost, testVirtualPath)
	assert.Contains(t, got.Auth, testAccessKey)
}

func TestConnect_Routing_AWSBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.us-east-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "mybucket.s3.us-east-1.amazonaws.com", testVirtualPath)
	assert.NotContains(t, got.Host, "mybucket.mybucket.")
}

func TestConnect_Routing_GCSPathStyle(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = gcsServiceEndpoint
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "storage.googleapis.com", testPathStylePath)
	assert.Equal(t, "https", got.Scheme)
}

func TestConnect_Routing_GCSBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = gcsServiceEndpoint + "/" + testBucket
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "storage.googleapis.com", testPathStylePath)
	assert.NotEqual(t, "/mybucket/mybucket/k", got.Path)
}

func TestConnect_Routing_MinioPathStyle(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://minio.internal.net"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, "minio.internal.net", testPathStylePath)
}

func TestConnect_Routing_LocalhostPathStyle(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "localhost"
	cfg.Port = testMinioPort
	cfg.Insecure = true
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assertSignedRequest(t, got, testMinioHost, testPathStylePath)
	assert.Equal(t, "http", got.Scheme)
}

func TestConnect_Routing_CustomS3LabelEndpoint(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://s3.minio.internal.net"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "s3.minio.internal.net", got.Host)
	assert.Equal(t, testPathStylePath, got.Path)
	assert.Contains(t, got.Auth, sigV4Algorithm)
}

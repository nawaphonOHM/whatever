package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_Routing_LegacyDashPathStyleBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://s3-us-west-2.amazonaws.com/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Contains(t, got.Host, "mybucket.s3")
	assert.Contains(t, got.Auth, "us-west-2")
}

func TestConnect_Routing_VPCEBucketURLPath(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://bucket.vpce-1a2b3c4d-5e6f.s3.us-east-1.vpce.amazonaws.com/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.bucket.vpce-1a2b3c4d-5e6f.s3.us-east-1.vpce.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "us-east-1")
}

func TestConnect_Routing_StroageTypoBucketURL(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://stroage.googleapis.com/mybucket"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "stroage.googleapis.com", got.Host)
	assert.Equal(t, "/mybucket/k", got.Path)
}

func TestConnect_Routing_ChinaVPCE(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.Endpoint = "https://bucket.vpce-1.s3.cn-north-1.vpce.amazonaws.com.cn"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.bucket.vpce-1.s3.cn-north-1.vpce.amazonaws.com.cn", got.Host)
	assert.Contains(t, got.Auth, "cn-north-1")
}

func TestConnect_Routing_Dualstack(t *testing.T) {
	t.Parallel()
	cfg := testConfigWithKeys()
	cfg.BucketURL = "https://mybucket.s3.dualstack.us-east-1.amazonaws.com"
	var got capturedRequest
	headWithConfig(t, cfg, &got)
	assert.Equal(t, "mybucket.s3.dualstack.us-east-1.amazonaws.com", got.Host)
	assert.Contains(t, got.Auth, "us-east-1")
}

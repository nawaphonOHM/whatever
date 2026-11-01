package client_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func assertNoChecksumHeaders(t *testing.T, req recordedRequest) {
	t.Helper()
	assert.Empty(t, req.Header.Get("x-amz-sdk-checksum-algorithm"))
	assert.Empty(t, req.Header.Get("x-amz-trailer"))
	assert.Empty(t, req.Header.Get("x-amz-checksum-crc32"))
	assert.Empty(t, req.Header.Get("x-amz-checksum-crc32c"))
	assert.Empty(t, req.Header.Get("x-amz-checksum-sha1"))
	assert.Empty(t, req.Header.Get("x-amz-checksum-sha256"))
	assert.NotContains(t, req.Header.Get("Content-Encoding"), "aws-chunked")
	assert.NotContains(t, req.Header.Get("X-Amz-Content-Sha256"), "TRAILER")
	assert.Empty(t, req.Header.Get("X-Amz-Decoded-Content-Length"))
}

func executePut(t *testing.T, cfg *config.Config) recordedRequest {
	t.Helper()
	rec := &recorder{respond: putOKResponder}
	c := connectRecorded(t, cfg, rec)
	_, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key:  testObjectKey,
		Body: strings.NewReader("sample payload data"),
	})
	require.NoError(t, err)
	require.Len(t, rec.all(), 1)
	return rec.all()[0]
}

func TestPutObject_GCS_WireNoUnsolicitedChecksums(t *testing.T) {
	t.Parallel()

	cfgEp := testConfigWithKeys()
	cfgEp.Endpoint = gcsServiceEndpoint
	assertNoChecksumHeaders(t, executePut(t, cfgEp))

	cfgURL := testConfigWithKeys()
	cfgURL.BucketURL = "https://mybucket.storage.googleapis.com"
	assertNoChecksumHeaders(t, executePut(t, cfgURL))
}

func TestPutObject_AWS_WireAddsChecksumTrailerContrast(t *testing.T) {
	setAWSEnvCredentials(t)

	cfg := config.DefaultConfig()
	cfg.Bucket = testBucket
	cfg.Region = "us-east-1"
	req := executePut(t, cfg)
	assert.Contains(t, req.Header.Get("Content-Encoding"), "aws-chunked")
	assert.Equal(t, "x-amz-checksum-crc32", req.Header.Get("x-amz-trailer"))
	assert.Equal(t, "STREAMING-UNSIGNED-PAYLOAD-TRAILER", req.Header.Get("X-Amz-Content-Sha256"))
}

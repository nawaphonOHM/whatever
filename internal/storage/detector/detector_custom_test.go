package detector

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDetectProvider_S3Compatible verifies detection of generic S3-compatible endpoints
// and ensures query, path, userinfo, and unrelated label substrings do not cause false positives.
func TestDetectProvider_S3Compatible(t *testing.T) {
	endpoints := []string{
		"http://localhost:9000",
		"http://localhost:9000/s3-backups",
		"https://minio.internal.net",
		"https://minio.internal.net/mys3bucket",
		"http://minio:9000?x=s3",
		"https://user:s3pass@minio.local",
		"https://minio.example.com/storage.googleapis",
		"https://mys3server.corp",
		"http://127.0.0.1:4566",
		"http://my-storage-server:8080",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			p, err := DetectProvider(ep, "")
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderS3Compatible, p)
		})
	}
}

// TestDetectProvider_S3Compatible_BucketURL verifies custom provider from bucketURL.
func TestDetectProvider_S3Compatible_BucketURL(t *testing.T) {
	p, err := DetectProvider("", "http://minio.internal.net/mybucket")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderS3Compatible, p)
}

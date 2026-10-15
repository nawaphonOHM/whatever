package detector

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDetectProvider_GCS verifies detection of Google Cloud Storage endpoints.
func TestDetectProvider_GCS(t *testing.T) {
	endpoints := []string{
		"https://storage.googleapis.com",
		"https://mybucket.storage.googleapis.com",
		"stroage.googleapis",
		"https://s3-backup.storage.googleapis.com",
		"https://storage.googleapis.com/mybucket/s3-data",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			p, err := DetectProvider(ep, "")
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderGCS, p)
		})
	}
}

// TestDetectProvider_GCS_FromBucketURL verifies GCS detection from bucketURL.
func TestDetectProvider_GCS_FromBucketURL(t *testing.T) {
	p, err := DetectProvider("", "https://mybucket.stroage.googleapis.com")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderGCS, p)
}

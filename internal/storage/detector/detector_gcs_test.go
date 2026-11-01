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
		"https://commondatastorage.googleapis.com",
		"https://mybucket.commondatastorage.googleapis.com",
		"https://storage.cloud.google.com",
		"https://storage.us-central1.rep.googleapis.com",
		"https://mybucket.storage.us-central1.rep.googleapis.com",
		"https://storage.p.googleapis.com",
		"https://storage-custom.p.googleapis.com",
		"https://mybucket.storage-psc.p.googleapis.com",
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

func TestBucketSubdomain_GCS(t *testing.T) {
	tests := []struct {
		host       string
		wantBucket string
		wantOK     bool
	}{
		{"storage.googleapis.com", "", false},
		{"mybucket.storage.googleapis.com", exampleBucket, true},
		{"storage-logs.storage.googleapis.com", "storage-logs", true},
		{"commondatastorage.googleapis.com", "", false},
		{"mybucket.commondatastorage.googleapis.com", exampleBucket, true},
		{"storage.us-central1.rep.googleapis.com", "", false},
		{"mybucket.storage.us-central1.rep.googleapis.com", exampleBucket, true},
		{"storage-logs.storage.us-central1.rep.googleapis.com", "storage-logs", true},
		{"storage-psc.p.googleapis.com", "", false},
		{"storage-custom.p.googleapis.com", "", false},
		{"mybucket.storage-psc.p.googleapis.com", exampleBucket, true},
		{"mybucket.storage-custom.p.googleapis.com", exampleBucket, true},
		{"storage-logs.storage-psc.p.googleapis.com", "storage-logs", true},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			gotBucket, gotOK := BucketSubdomain(tt.host)
			assert.Equal(t, tt.wantOK, gotOK)
			assert.Equal(t, tt.wantBucket, gotBucket)
		})
	}
}

package client_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_GCSBucketMismatch(t *testing.T) {
	t.Parallel()
	runMismatchCases(t, contracts.ProviderGCS, endpointCfg, []mismatchCase{
		{"storage-logs.storage.googleapis.com", "storage-logs"},
		{"Storage-Logs.Storage.GoogleAPIs.com", "storage-logs"},
		{"storage-logs.storage.googleapis.com.", "storage-logs"},
		{"other.storage.googleapis.com.", otherBucket},
		{"other.storage-psc.p.googleapis.com", otherBucket},
		{"other.storage-psc.p.googleapis.com.", otherBucket},
		{"other.storage.us-central1.rep.googleapis.com", otherBucket},
		{"other.storage.us-central1.rep.googleapis.com.", otherBucket},
	})
}

func TestBuildS3Options_GCSBucketURLMismatch(t *testing.T) {
	t.Parallel()
	runMismatchCases(t, contracts.ProviderGCS, bucketURLCfg, []mismatchCase{
		{"other.storage.googleapis.com.", otherBucket},
		{"other.storage-psc.p.googleapis.com/mybucket", otherBucket},
	})
}

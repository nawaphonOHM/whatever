package client_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_GCSEndpointRouting(t *testing.T) {
	t.Parallel()
	runRoutingCases(t, contracts.ProviderGCS, endpointCfg, []routingCase{
		{"MyBucket.Storage.GoogleAPIs.com", "https://Storage.GoogleAPIs.com", "", true},
		{"MYBUCKET.STORAGE-PSC.P.GoogleAPIs.com", "https://STORAGE-PSC.P.GoogleAPIs.com", "", true},
		{"MyBucket.Storage.US-CENTRAL1.REP.GoogleAPIs.com", "https://Storage.US-CENTRAL1.REP.GoogleAPIs.com", "", true},
		{"mybucket.storage.us-central1.rep.googleapis.com", gcsRegionalBase, "", true},
		{"storage.us-central1.rep.googleapis.com", gcsRegionalBase, "", true},
		{"mybucket.storage-psc.p.googleapis.com", gcsPSCBase, "", true},
		{"storage-psc.p.googleapis.com", gcsPSCBase, "", true},
		{"mybucket.storage-custom.p.googleapis.com", gcsCustomBase, "", true},
		{"storage-custom.p.googleapis.com", gcsCustomBase, "", true},
	})
}

func TestBuildS3Options_GCSBucketURLRouting(t *testing.T) {
	t.Parallel()
	runRoutingCases(t, contracts.ProviderGCS, bucketURLCfg, []routingCase{
		{"MyBucket.Storage.GoogleAPIs.com", "https://Storage.GoogleAPIs.com", "", true},
		{"MYBUCKET.STORAGE-PSC.P.GoogleAPIs.com", "https://STORAGE-PSC.P.GoogleAPIs.com", "", true},
		{"MyBucket.Storage.US-CENTRAL1.REP.GoogleAPIs.com", "https://Storage.US-CENTRAL1.REP.GoogleAPIs.com", "", true},
		{"mybucket.storage-psc.p.googleapis.com", gcsPSCBase, "", true},
		{"storage-psc.p.googleapis.com/mybucket", gcsPSCBase, "", true},
	})
}

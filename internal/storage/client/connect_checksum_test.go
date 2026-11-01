package client_test

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func TestConnect_Checksum_GCS_Static(t *testing.T) {
	t.Parallel()

	cfgEp := testConfigWithKeys()
	cfgEp.Endpoint = gcsServiceEndpoint
	assertConnectChecksumWhenRequired(t, cfgEp)

	cfgURL := testConfigWithKeys()
	cfgURL.BucketURL = "https://mybucket.storage.googleapis.com"
	assertConnectChecksumWhenRequired(t, cfgURL)
}

func TestConnect_Checksum_GCS_DefaultChain(t *testing.T) {
	setAWSEnvCredentials(t)

	cfgEp := config.DefaultConfig()
	cfgEp.Bucket = testBucket
	cfgEp.Endpoint = gcsServiceEndpoint
	assertConnectChecksumWhenRequired(t, cfgEp)

	cfgURL := config.DefaultConfig()
	cfgURL.Bucket = testBucket
	cfgURL.BucketURL = "https://mybucket.storage.googleapis.com"
	assertConnectChecksumWhenRequired(t, cfgURL)
}

func TestConnect_Checksum_GCS_EnvOverridesIgnored(t *testing.T) {
	setOpposingChecksumEnv(t)

	cfgStatic := testConfigWithKeys()
	cfgStatic.Endpoint = gcsServiceEndpoint
	assertConnectChecksumWhenRequired(t, cfgStatic)

	cfgDefault := config.DefaultConfig()
	cfgDefault.Bucket = testBucket
	cfgDefault.Endpoint = gcsServiceEndpoint
	assertConnectChecksumWhenRequired(t, cfgDefault)
}

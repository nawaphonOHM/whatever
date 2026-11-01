package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func malformedSecurityCases() []struct {
	name      string
	endpoint  string
	bucketURL string
} {
	return []struct {
		name      string
		endpoint  string
		bucketURL string
	}{
		{
			name:      "whitespace in endpoint",
			endpoint:  secAdminPrefix + secPassword + "@minio .internal:9000/data?" + secToken,
			bucketURL: "https://" + secUserPrefix + secPassword + "@storage.googleapis.com",
		},
		{
			name:      "tab in bucket URL",
			endpoint:  "",
			bucketURL: secUserPrefix + secPassword + "@storage.googleapis\t.com?" + secToken,
		},
		{
			name:      "schemeless with only credentials",
			endpoint:  secUserPrefix + secPassword + secAt,
			bucketURL: "https://s3.amazonaws.com",
		},
		{
			name:      "percent encoding error in credentials",
			endpoint:  "https://%zz:" + secPassword + secAt + secHostMinio + ":9000",
			bucketURL: "",
		},
		{
			name:      "single slash https with credentials in endpoint",
			endpoint:  "https:/" + "admin:" + secPassword + secAt + secHostMinio + ":9000",
			bucketURL: "",
		},
		{
			name:      "single slash http with credentials and port",
			endpoint:  "http:/admin:" + secPassword + "@minio:9000",
			bucketURL: "",
		},
		{
			name:      "single slash s3 with credentials in bucket url",
			endpoint:  "",
			bucketURL: "s3:/" + secUserPrefix + secPassword + secAt + "host/x",
		},
		{
			name:      "secret in port position with absolute endpoint and credentials",
			endpoint:  secAdminPrefix + secPassword + secAt + secHostMinio + ":" + secPassword,
			bucketURL: "",
		},
		{
			name:      "secret in port position with absolute endpoint without credentials",
			endpoint:  "https://host:" + secPassword,
			bucketURL: "",
		},
		{
			name:      "secret in port position with schemeless endpoint",
			endpoint:  secHostMinio + ":" + secPassword,
			bucketURL: "",
		},
	}
}

func TestConnect_MalformedSecurityInputs_RejectionsDoNotLeakSecrets(t *testing.T) {
	for _, tc := range malformedSecurityCases() {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = secTestBucket
			cfg.Region = secTestRegion
			cfg.AccessKeyID = "test-key"
			cfg.SecretAccessKey = secSecretKey
			cfg.Endpoint = tc.endpoint
			cfg.BucketURL = tc.bucketURL

			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			assertMalformedRejection(t, cfg, buf.String(), err)
		})
	}
}

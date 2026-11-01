package client

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func TestConnect_AcceptedInputsWithSecretPayload_DiagnosticLogsSanitized(t *testing.T) {
	testCases := []struct {
		name              string
		endpoint          string
		bucketURL         string
		bucketName        string
		expectRawEndpoint string
		expectRawBucket   string
		expectHost        string
		expectProvider    string
	}{
		{
			name: "absolute endpoint and bucket URL with credentials and query",
			endpoint: secAdminPrefix + secPassword + "@" + secHostMinio +
				":9000/data?token=" + secToken + secHashFrag,
			bucketURL: "https://" + secUserPrefix + secPassword +
				"@bucket.storage.googleapis.com?token=" + secToken + secHashFrag,
			expectRawEndpoint: "https://minio.internal:9000/data",
			expectRawBucket:   "https://bucket.storage.googleapis.com",
			expectHost:        secHostMinio,
			expectProvider:    secProviderS3,
		},
		{
			name: "schemeless endpoint with query scheme separator and nested redirect",
			endpoint: "admin:" + secPassword + "@" + secHostMinio +
				":9000/data?redir=https://other.internal/x&token=" + secToken + secHashFrag,
			bucketURL: secUserPrefix + secPassword +
				"@custom-host:9000/data?r=a://b&token=" + secToken + "#frag://b",
			expectRawEndpoint: "minio.internal:9000/data",
			expectRawBucket:   "custom-host:9000/data",
			expectHost:        secHostMinio,
			expectProvider:    secProviderS3,
		},
		{
			name:     "schemeless aws bucket URL with credentials",
			endpoint: "",
			bucketURL: secUserPrefix + secPassword +
				"@mybucket.s3.us-east-1.amazonaws.com/data?token=" + secToken,
			bucketName:        "mybucket",
			expectRawEndpoint: "",
			expectRawBucket:   "mybucket.s3.us-east-1.amazonaws.com/data",
			expectHost:        "mybucket.s3.us-east-1.amazonaws.com",
			expectProvider:    "aws",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWSEnv(t)
			var buf bytes.Buffer
			setStorageTestLogger(t, storageLogModes()[0].debugHandler(&buf))

			cfg := config.DefaultConfig()
			cfg.Bucket = secTestBucket
			if tc.bucketName != "" {
				cfg.Bucket = tc.bucketName
			}
			cfg.Region = secTestRegion
			cfg.AccessKeyID = secTestAccessKey
			cfg.SecretAccessKey = secSecretKey
			cfg.SessionToken = secSessionToken
			cfg.Endpoint = tc.endpoint
			cfg.BucketURL = tc.bucketURL

			_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
			require.NoError(t, err)

			assertNoSensitiveDataInOutput(t, buf.String())

			entries := parseStorageLogs(t, buf.Bytes())
			detectLog := findStorageLog(entries, secLogEvalDetect)
			require.NotNil(t, detectLog, "diagnostic detection log must be emitted for accepted inputs")
			assert.Equal(t, tc.expectRawEndpoint, detectLog["raw_endpoint"])
			assert.Equal(t, tc.expectRawBucket, detectLog["raw_bucket_url"])
			assert.Equal(t, tc.expectHost, detectLog[secLogInspectHost])
			assert.Equal(t, tc.expectProvider, detectLog[secLogDetectProv])
		})
	}
}

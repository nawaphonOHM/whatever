package detector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

func TestDetect_ComprehensiveCases(t *testing.T) {
	cases := []struct {
		name         string
		endpoint     string
		bucketURL    string
		wantProvider contracts.Provider
		wantRule     string
		wantHost     string
		wantBucket   string
		wantErr      bool
	}{
		// Default rule cases
		{
			name:         "empty endpoint and empty bucketURL defaults to AWS",
			endpoint:     "",
			bucketURL:    "",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleDefaultAWS,
			wantHost:     "",
			wantBucket:   "",
		},

		// AWS S3 standard and regional cases
		{
			name:         "AWS global endpoint HTTPS",
			endpoint:     "https://s3.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS regional endpoint us-east-1",
			endpoint:     "https://s3.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.us-east-1.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS regional legacy hyphen endpoint s3-us-west-2",
			endpoint:     "https://s3-us-west-2.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3-us-west-2.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS schemeless regional endpoint",
			endpoint:     "s3.ap-southeast-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.ap-southeast-1.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS virtual-host bucket endpoint",
			endpoint:     "https://mybucket.s3.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "mybucket.s3.us-east-1.amazonaws.com",
			wantBucket:   "mybucket",
		},
		{
			name:         "AWS virtual-host bucket with uppercase and trailing dot",
			endpoint:     "https://MyBucket.S3.us-east-1.amazonaws.com.",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "mybucket.s3.us-east-1.amazonaws.com",
			wantBucket:   "mybucket",
		},
		{
			name:         "AWS dualstack endpoint",
			endpoint:     "https://s3.dualstack.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.dualstack.us-east-1.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS api.aws domain endpoint",
			endpoint:     "https://s3.dualstack.us-east-1.api.aws",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.dualstack.us-east-1.api.aws",
			wantBucket:   "",
		},
		{
			name:         "AWS China region amazonaws.com.cn",
			endpoint:     "https://s3.cn-north-1.amazonaws.com.cn",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3.cn-north-1.amazonaws.com.cn",
			wantBucket:   "",
		},
		{
			name:         "AWS S3 Express One Zone base endpoint",
			endpoint:     "https://s3express-use1-az4.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "s3express-use1-az4.us-east-1.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS S3 Express One Zone directory bucket endpoint",
			endpoint:     "https://sample--use1-az4--x-s3.s3express-use1-az4.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "sample--use1-az4--x-s3.s3express-use1-az4.us-east-1.amazonaws.com",
			wantBucket:   "sample--use1-az4--x-s3",
		},
		{
			name:         "AWS PrivateLink VPCE base endpoint",
			endpoint:     "https://bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "AWS PrivateLink VPCE bucket endpoint",
			endpoint:     "https://mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantProvider: contracts.ProviderAWS,
			wantRule:     detector.RuleAWSHostname,
			wantHost:     "mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantBucket:   "mybucket",
		},

		// Google Cloud Storage cases
		{
			name:         "GCS standard domain",
			endpoint:     "https://storage.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "storage.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "GCS schemeless standard domain",
			endpoint:     "storage.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "storage.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "GCS virtual-hosted bucket domain",
			endpoint:     "https://mybucket.storage.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "mybucket.storage.googleapis.com",
			wantBucket:   "mybucket",
		},
		{
			name:         "GCS regional rep domain",
			endpoint:     "https://storage.us-central1.rep.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "storage.us-central1.rep.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "GCS regional rep bucket domain",
			endpoint:     "https://mybucket.storage.us-central1.rep.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "mybucket.storage.us-central1.rep.googleapis.com",
			wantBucket:   "mybucket",
		},
		{
			name:         "GCS PSC base domain",
			endpoint:     "https://storage-psc.p.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "storage-psc.p.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "GCS PSC bucket domain",
			endpoint:     "https://mybucket.storage-psc.p.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "mybucket.storage-psc.p.googleapis.com",
			wantBucket:   "mybucket",
		},
		{
			name:         "GCS commondatastorage domain",
			endpoint:     "https://commondatastorage.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "commondatastorage.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "GCS cloud google com domain",
			endpoint:     "https://storage.cloud.google.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     detector.RuleGCSHostname,
			wantHost:     "storage.cloud.google.com",
			wantBucket:   "",
		},

		// Generic S3-compatible fallback cases
		{
			name:         "Local MinIO with port",
			endpoint:     "http://localhost:9000",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "localhost",
			wantBucket:   "",
		},
		{
			name:         "IP address endpoint",
			endpoint:     "http://127.0.0.1:9000",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "127.0.0.1",
			wantBucket:   "",
		},
		{
			name:         "MinIO domain with path prefix",
			endpoint:     "https://minio.internal.net/s3",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "minio.internal.net",
			wantBucket:   "",
		},
		{
			name:         "Custom proxy domain",
			endpoint:     "https://s3-gateway.example.com",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "s3-gateway.example.com",
			wantBucket:   "",
		},
		{
			name:         "Non-S3 AWS endpoint",
			endpoint:     "https://dynamodb.us-east-1.amazonaws.com",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "dynamodb.us-east-1.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "Non-GCS GCP endpoint",
			endpoint:     "https://compute.googleapis.com",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     detector.RuleS3CompatibleFallback,
			wantHost:     "compute.googleapis.com",
			wantBucket:   "",
		},

		// Error cases
		{
			name:     "Invalid URL with whitespace",
			endpoint: "http://localhost:9000 /path",
			wantErr:  true,
		},
		{
			name:     "Invalid port out of range",
			endpoint: "http://localhost:70000",
			wantErr:  true,
		},
		{
			name:     "Unsupported scheme ftp",
			endpoint: "ftp://storage.example.com",
			wantErr:  true,
		},
		{
			name:     "Missing host in URL",
			endpoint: "http://",
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := detector.Detect(tc.endpoint, tc.bucketURL)
			if tc.wantErr {
				require.Error(t, err)
				assert.Empty(t, d.Provider)
				assert.Empty(t, d.Rule)

				p, err := detector.DetectProvider(tc.endpoint, tc.bucketURL)
				require.Error(t, err)
				assert.Empty(t, p)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantProvider, d.Provider)
			assert.Equal(t, tc.wantRule, d.Rule)
			assert.Equal(t, tc.wantHost, d.Host)
			assert.Equal(t, tc.wantBucket, d.Bucket)

			p, err := detector.DetectProvider(tc.endpoint, tc.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, tc.wantProvider, p)
		})
	}
}

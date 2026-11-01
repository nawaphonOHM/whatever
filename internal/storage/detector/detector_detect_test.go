package detector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/detector"
)

func TestDetect_DefaultAWS(t *testing.T) {
	d, err := detector.Detect("", "")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, d.Provider)
	assert.Equal(t, detector.RuleDefaultAWS, d.Rule)
	assert.Empty(t, d.Host)
	assert.Empty(t, d.Bucket)

	p, err := detector.DetectProvider("", "")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, p)
}

func TestDetect_AWSHostnames(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		bucketURL  string
		wantHost   string
		wantBucket string
	}{
		{
			name:       "standard global AWS S3 endpoint",
			endpoint:   "https://s3.amazonaws.com",
			wantHost:   "s3.amazonaws.com",
			wantBucket: "",
		},
		{
			name:       "regional AWS S3 endpoint",
			endpoint:   "https://s3.us-west-2.amazonaws.com",
			wantHost:   "s3.us-west-2.amazonaws.com",
			wantBucket: "",
		},
		{
			name:       "virtual-hosted bucket AWS S3 endpoint",
			endpoint:   "https://test-bucket.s3.us-east-1.amazonaws.com",
			wantHost:   "test-bucket.s3.us-east-1.amazonaws.com",
			wantBucket: "test-bucket",
		},
		{
			name:       "AWS S3 Express One Zone endpoint",
			endpoint:   "https://mybucket--use1-az4--x-s3.s3express-use1-az4.us-east-1.amazonaws.com",
			wantHost:   "mybucket--use1-az4--x-s3.s3express-use1-az4.us-east-1.amazonaws.com",
			wantBucket: "mybucket--use1-az4--x-s3",
		},
		{
			name:       "AWS China region S3 endpoint",
			endpoint:   "https://data-bucket.s3.cn-north-1.amazonaws.com.cn",
			wantHost:   "data-bucket.s3.cn-north-1.amazonaws.com.cn",
			wantBucket: "data-bucket",
		},
		{
			name:       "AWS FIPS endpoint",
			endpoint:   "https://s3-fips.us-east-1.amazonaws.com",
			wantHost:   "s3-fips.us-east-1.amazonaws.com",
			wantBucket: "",
		},
		{
			name:       "AWS PrivateLink VPCE endpoint with bucket prefix",
			endpoint:   "https://mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantHost:   "mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
			wantBucket: "mybucket",
		},
		{
			name:       "AWS detection from bucketURL fallback",
			endpoint:   "",
			bucketURL:  "https://archive.s3.ap-southeast-1.amazonaws.com",
			wantHost:   "archive.s3.ap-southeast-1.amazonaws.com",
			wantBucket: "archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := detector.Detect(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderAWS, d.Provider)
			assert.Equal(t, detector.RuleAWSHostname, d.Rule)
			assert.Equal(t, tt.wantHost, d.Host)
			assert.Equal(t, tt.wantBucket, d.Bucket)

			p, err := detector.DetectProvider(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderAWS, p)
		})
	}
}

func TestDetect_GCSHostnames(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		bucketURL  string
		wantHost   string
		wantBucket string
	}{
		{
			name:       "standard GCS endpoint",
			endpoint:   "https://storage.googleapis.com",
			wantHost:   "storage.googleapis.com",
			wantBucket: "",
		},
		{
			name:       "virtual-hosted bucket GCS endpoint",
			endpoint:   "https://my-bucket.storage.googleapis.com",
			wantHost:   "my-bucket.storage.googleapis.com",
			wantBucket: "my-bucket",
		},
		{
			name:       "GCS regional rep endpoint",
			endpoint:   "https://assets.storage.us-central1.rep.googleapis.com",
			wantHost:   "assets.storage.us-central1.rep.googleapis.com",
			wantBucket: "assets",
		},
		{
			name:       "GCS private service connect PSC endpoint",
			endpoint:   "https://media.storage-psc.p.googleapis.com",
			wantHost:   "media.storage-psc.p.googleapis.com",
			wantBucket: "media",
		},
		{
			name:       "GCS commondatastorage endpoint",
			endpoint:   "https://commondatastorage.googleapis.com",
			wantHost:   "commondatastorage.googleapis.com",
			wantBucket: "",
		},
		{
			name:       "GCS detection from bucketURL fallback",
			endpoint:   "",
			bucketURL:  "https://my-bucket.storage.googleapis.com",
			wantHost:   "my-bucket.storage.googleapis.com",
			wantBucket: "my-bucket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := detector.Detect(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderGCS, d.Provider)
			assert.Equal(t, detector.RuleGCSHostname, d.Rule)
			assert.Equal(t, tt.wantHost, d.Host)
			assert.Equal(t, tt.wantBucket, d.Bucket)

			p, err := detector.DetectProvider(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderGCS, p)
		})
	}
}

func TestDetect_S3CompatibleFallback(t *testing.T) {
	tests := []struct {
		name      string
		endpoint  string
		bucketURL string
		wantHost  string
	}{
		{
			name:     "local MinIO instance with port",
			endpoint: "http://localhost:9000",
			wantHost: "localhost",
		},
		{
			name:     "internal MinIO DNS",
			endpoint: "https://minio.corp.internal",
			wantHost: "minio.corp.internal",
		},
		{
			name:     "LocalStack simulated S3",
			endpoint: "http://127.0.0.1:4566",
			wantHost: "127.0.0.1",
		},
		{
			name:     "Wasabi cloud storage",
			endpoint: "https://s3.wasabisys.com",
			wantHost: "s3.wasabisys.com",
		},
		{
			name:     "Backblaze B2 S3 API",
			endpoint: "https://s3.us-west-004.backblazeb2.com",
			wantHost: "s3.us-west-004.backblazeb2.com",
		},
		{
			name:     "Custom reverse proxy for S3",
			endpoint: "https://storage-proxy.mycompany.org",
			wantHost: "storage-proxy.mycompany.org",
		},
		{
			name:     "AWS EC2 endpoint is not S3",
			endpoint: "https://ec2.us-east-1.amazonaws.com",
			wantHost: "ec2.us-east-1.amazonaws.com",
		},
		{
			name:     "GCP BigQuery endpoint is not GCS",
			endpoint: "https://bigquery.googleapis.com",
			wantHost: "bigquery.googleapis.com",
		},
		{
			name:      "S3 compatible fallback from bucketURL",
			endpoint:  "",
			bucketURL: "http://minio.local:9000/mybucket",
			wantHost:  "minio.local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := detector.Detect(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderS3Compatible, d.Provider)
			assert.Equal(t, detector.RuleS3CompatibleFallback, d.Rule)
			assert.Equal(t, tt.wantHost, d.Host)
			assert.Empty(t, d.Bucket)

			p, err := detector.DetectProvider(tt.endpoint, tt.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderS3Compatible, p)
		})
	}
}

func TestDetect_Errors(t *testing.T) {
	badTargets := []string{
		"http://invalid host with spaces.com",
		"http://",
		"ftp://unsupported-scheme.com",
		"http://localhost:99999",
		"http://localhost:0",
		"http://localhost:-5",
	}

	for _, bad := range badTargets {
		t.Run(bad, func(t *testing.T) {
			d, err := detector.Detect(bad, "")
			require.Error(t, err)
			assert.Empty(t, d.Provider)
			assert.Empty(t, d.Rule)

			p, err := detector.DetectProvider(bad, "")
			require.Error(t, err)
			assert.Empty(t, p)
		})
	}
}

func TestBucketSubdomain_GenericS3CompatibleNeverExtractsBucket(t *testing.T) {
	customHosts := []string{
		"mybucket.localhost",
		"mybucket.minio.internal",
		"mybucket.s3.wasabisys.com",
		"mybucket.custom-domain.com",
		"bucket.localhost:9000",
	}

	for _, host := range customHosts {
		t.Run(host, func(t *testing.T) {
			bucket, ok := detector.BucketSubdomain(host)
			assert.False(t, ok)
			assert.Empty(t, bucket)
		})
	}
}

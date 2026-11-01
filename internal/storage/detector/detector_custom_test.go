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
		"https://s3.wasabisys.com",
		"https://s3.us-west-004.backblazeb2.com",
		"https://s3.fr-par.scw.cloud",
		"https://s3.us-south.cloud-object-storage.appdomain.cloud",
		"https://s3.minio.corp.local",
		"https://s3.example.com",
		"https://s3-custom.mycompany.org",
		"https://s3.files.example.com",
		"https://bigquery.googleapis.com",
		"https://compute.p.googleapis.com",
		"https://logging.rep.googleapis.com",
		"http://ec2-54-1-2-3.compute-1.amazonaws.com:9000",
		"https://minio-lb-123.us-east-1.elb.amazonaws.com",
		"https://abc.execute-api.us-east-1.amazonaws.com",
		"https://ec2.us-east-1.api.aws",
		"https://s3-gw-123456.us-east-1.elb.amazonaws.com",
		"https://s3-proxy.execute-api.us-east-1.amazonaws.com",
		"https://s3.minio-lb.us-east-1.elb.amazonaws.com",
		"https://s3.us-east-1.elb.amazonaws.com",
		"https://s3-gw.compute-1.amazonaws.com",
		"https://amazonaws.com",
		"https://bigquerystorage.us-central1.rep.googleapis.com",
		"https://bigquerystorage-ep.p.googleapis.com",
		"https://storagetransfer.us-central1.rep.googleapis.com",
		"https://my-storage.compute-ep.p.googleapis.com",
		"https://my-storage.compute.us-central1.rep.googleapis.com",
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

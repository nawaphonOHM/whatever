package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleBucket = "mybucket"

func TestHostClassification_NormalizesTrailingDotAndCase(t *testing.T) {
	gcsHosts := []string{
		"mybucket.storage.googleapis.com.",
		"MyBucket.Storage.GoogleAPIs.com.",
		"storage-psc.p.googleapis.com.",
		"storage.us-central1.rep.googleapis.com.",
	}
	for _, host := range gcsHosts {
		assert.True(t, IsGCSHost(host), host)
		assert.False(t, IsAWSHost(host), host)
	}

	awsHosts := []string{
		"mybucket.s3.us-east-1.amazonaws.com.",
		"MyBucket.S3.Amazonaws.com.",
		"s3express-use1-az4.us-east-1.amazonaws.com.",
	}
	for _, host := range awsHosts {
		assert.True(t, IsAWSHost(host), host)
		assert.False(t, IsGCSHost(host), host)
	}
}

func TestBucketSubdomain_NormalizesTrailingDot(t *testing.T) {
	tests := []struct {
		host       string
		wantBucket string
		wantOK     bool
	}{
		{"storage-logs.storage.googleapis.com.", "storage-logs", true},
		{"Other.Storage-PSC.P.GoogleAPIs.com.", "other", true},
		{"storage-psc.p.googleapis.com.", "", false},
		{"other.s3.us-east-1.amazonaws.com.", "other", true},
		{"s3.us-east-1.amazonaws.com.", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			gotBucket, gotOK := BucketSubdomain(tt.host)
			assert.Equal(t, tt.wantOK, gotOK)
			assert.Equal(t, tt.wantBucket, gotBucket)
		})
	}
}

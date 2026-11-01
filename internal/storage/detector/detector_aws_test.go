package detector

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDetectProvider_AWS verifies detection of AWS S3 endpoints.
func TestDetectProvider_AWS(t *testing.T) {
	endpoints := []string{
		"https://s3.amazonaws.com",
		"https://mybucket.s3.us-east-1.amazonaws.com",
		"s3.ap-southeast-1.amazonaws.com",
		"https://mybucket.s3-us-west-2.amazonaws.com",
		"https://s3-fips.us-east-1.amazonaws.com",
		"s3.dualstack.us-east-1.amazonaws.com",
		"https://s3.cn-north-1.amazonaws.com.cn",
		"https://mybucket.s3.cn-northwest-1.amazonaws.com.cn",
		"https://s3.dualstack.us-east-1.api.aws",
		"https://s3-accesspoint.us-east-1.amazonaws.com",
		"https://s3-outposts.us-east-1.amazonaws.com",
		"https://s3-object-lambda.us-east-1.amazonaws.com",
		"https://s3-control.us-east-1.amazonaws.com",
		"https://s3-website.us-east-1.amazonaws.com",
		"https://s3express-use1-az4.us-east-1.amazonaws.com",
		"https://mybucket--use1-az4--x-s3.s3express-use1-az4.us-east-1.amazonaws.com",
		"https://bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			p, err := DetectProvider(ep, "")
			require.NoError(t, err)
			assert.Equal(t, contracts.ProviderAWS, p)
		})
	}
}

// TestDetectProvider_DefaultAWS verifies default to AWS when both are empty.
func TestDetectProvider_DefaultAWS(t *testing.T) {
	p, err := DetectProvider("", "")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, p)
}

// TestDetectProvider_AWS_FromBucketURL verifies AWS detection from bucketURL.
func TestDetectProvider_AWS_FromBucketURL(t *testing.T) {
	p, err := DetectProvider("", "https://mybucket.s3.us-west-2.amazonaws.com")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, p)
}

func TestBucketSubdomain_AWS(t *testing.T) {
	tests := []struct {
		host       string
		wantBucket string
		wantOK     bool
	}{
		{"s3.amazonaws.com", "", false},
		{"s3.us-east-1.amazonaws.com", "", false},
		{"s3-us-west-2.amazonaws.com", "", false},
		{"mybucket.s3.amazonaws.com", exampleBucket, true},
		{"mybucket.s3.us-east-1.amazonaws.com", exampleBucket, true},
		{"mybucket.s3-us-west-2.amazonaws.com", exampleBucket, true},
		{"s3express-use1-az4.us-east-1.amazonaws.com", "", false},
		{"mybucket.s3express-use1-az4.us-east-1.amazonaws.com", exampleBucket, true},
		{"s3express-x.s3.us-east-1.amazonaws.com", "s3express-x", true},
		{"vpce-logs.s3.us-east-1.amazonaws.com", "vpce-logs", true},
		{"vpce-logs.s3.amazonaws.com", "vpce-logs", true},
		{"vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", "", false},
		{"bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", "", false},
		{"bucket.vpce-1.s3.cn-north-1.vpce.amazonaws.com.cn", "", false},
		{"mybucket.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", exampleBucket, true},
		{"other.bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", "other", true},
		{"other.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", "other", true},
		{"dynamodb.us-east-1.amazonaws.com", "", false},
		{"mybucket.dynamodb.us-east-1.amazonaws.com", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			gotBucket, gotOK := BucketSubdomain(tt.host)
			assert.Equal(t, tt.wantOK, gotOK)
			assert.Equal(t, tt.wantBucket, gotBucket)
		})
	}
}

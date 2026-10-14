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

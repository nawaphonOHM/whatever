package detector

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDetectProvider_Precedence verifies endpoint parameter takes priority over bucketURL.
func TestDetectProvider_Precedence(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		bucketURL  string
		expected   contracts.Provider
		expectRule string
		expectHost string
	}{
		{
			name:       "endpoint GCS overrides bucketURL AWS",
			endpoint:   "https://storage.googleapis.com",
			bucketURL:  "https://s3.amazonaws.com",
			expected:   contracts.ProviderGCS,
			expectRule: RuleGCSHostname,
			expectHost: "storage.googleapis.com",
		},
		{
			name:       "endpoint AWS overrides bucketURL GCS",
			endpoint:   "https://s3.amazonaws.com",
			bucketURL:  "https://storage.googleapis.com",
			expected:   contracts.ProviderAWS,
			expectRule: RuleAWSHostname,
			expectHost: "s3.amazonaws.com",
		},
		{
			name:       "endpoint custom overrides bucketURL AWS",
			endpoint:   "http://localhost:9000",
			bucketURL:  "https://s3.amazonaws.com",
			expected:   contracts.ProviderS3Compatible,
			expectRule: RuleS3CompatibleFallback,
			expectHost: "localhost",
		},
		{
			name:       "endpoint AWS overrides bucketURL custom",
			endpoint:   "https://s3.us-west-2.amazonaws.com",
			bucketURL:  "http://minio.local:9000",
			expected:   contracts.ProviderAWS,
			expectRule: RuleAWSHostname,
			expectHost: "s3.us-west-2.amazonaws.com",
		},
		{
			name:       "endpoint custom overrides bucketURL GCS",
			endpoint:   "http://minio.local:9000",
			bucketURL:  "https://storage.googleapis.com",
			expected:   contracts.ProviderS3Compatible,
			expectRule: RuleS3CompatibleFallback,
			expectHost: "minio.local",
		},
		{
			name:       "empty endpoint uses bucketURL GCS",
			endpoint:   "",
			bucketURL:  "https://storage.googleapis.com",
			expected:   contracts.ProviderGCS,
			expectRule: RuleGCSHostname,
			expectHost: "storage.googleapis.com",
		},
		{
			name:       "empty endpoint uses bucketURL AWS",
			endpoint:   "",
			bucketURL:  "https://s3.eu-central-1.amazonaws.com",
			expected:   contracts.ProviderAWS,
			expectRule: RuleAWSHostname,
			expectHost: "s3.eu-central-1.amazonaws.com",
		},
		{
			name:       "empty endpoint uses bucketURL custom",
			endpoint:   "",
			bucketURL:  "http://custom-s3.corp:8080/bucket",
			expected:   contracts.ProviderS3Compatible,
			expectRule: RuleS3CompatibleFallback,
			expectHost: "custom-s3.corp",
		},
		{
			name:       "both empty defaults to AWS",
			endpoint:   "",
			bucketURL:  "",
			expected:   contracts.ProviderAWS,
			expectRule: RuleDefaultAWS,
			expectHost: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := DetectProvider(tc.endpoint, tc.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, p)

			d, err := Detect(tc.endpoint, tc.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, d.Provider)
			assert.Equal(t, tc.expectRule, d.Rule)
			assert.Equal(t, tc.expectHost, d.Host)
		})
	}
}

// TestDetect_PrecedenceIgnoresUnusedTarget verifies that malformed unused parameter is ignored.
func TestDetect_PrecedenceIgnoresUnusedTarget(t *testing.T) {
	// When endpoint is valid, malformed bucketURL is ignored
	d, err := Detect("https://s3.amazonaws.com", "http://bad host.com")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, d.Provider)
	assert.Equal(t, RuleAWSHostname, d.Rule)
	assert.Equal(t, "s3.amazonaws.com", d.Host)

	p, err := DetectProvider("https://s3.amazonaws.com", "http://bad host.com")
	require.NoError(t, err)
	assert.Equal(t, contracts.ProviderAWS, p)
}

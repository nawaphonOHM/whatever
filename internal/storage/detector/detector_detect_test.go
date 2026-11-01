package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestDetect_RulesAndProviders(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		bucketURL    string
		wantProvider contracts.Provider
		wantRule     string
		wantHost     string
		wantBucket   string
	}{
		{
			name:         "empty inputs default to AWS",
			endpoint:     "",
			bucketURL:    "",
			wantProvider: contracts.ProviderAWS,
			wantRule:     RuleDefaultAWS,
			wantHost:     "",
			wantBucket:   "",
		},
		{
			name:         "gcs base endpoint",
			endpoint:     "https://storage.googleapis.com",
			bucketURL:    "",
			wantProvider: contracts.ProviderGCS,
			wantRule:     RuleGCSHostname,
			wantHost:     "storage.googleapis.com",
			wantBucket:   "",
		},
		{
			name:         "gcs virtual-hosted bucket URL",
			endpoint:     "",
			bucketURL:    "https://sample-bucket.storage.googleapis.com",
			wantProvider: contracts.ProviderGCS,
			wantRule:     RuleGCSHostname,
			wantHost:     "sample-bucket.storage.googleapis.com",
			wantBucket:   "sample-bucket",
		},
		{
			name:         "aws base endpoint",
			endpoint:     "https://s3.us-west-2.amazonaws.com",
			bucketURL:    "",
			wantProvider: contracts.ProviderAWS,
			wantRule:     RuleAWSHostname,
			wantHost:     "s3.us-west-2.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "s3-compatible custom endpoint",
			endpoint:     "https://custom-storage.internal.net",
			bucketURL:    "",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     RuleS3CompatibleFallback,
			wantHost:     "custom-storage.internal.net",
			wantBucket:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := Detect(tc.endpoint, tc.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, tc.wantProvider, d.Provider)
			assert.Equal(t, tc.wantRule, d.Rule)
			assert.Equal(t, tc.wantHost, d.Host)
			assert.Equal(t, tc.wantBucket, d.Bucket)
		})
	}
}

package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestDetect_RulesAndProviders_SchemelessWithQuery(t *testing.T) {
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
			name:         "schemeless with query containing scheme separator and credentials",
			endpoint:     "user:secretpass@custom-storage.internal.net:9000/data?r=a://b",
			bucketURL:    "",
			wantProvider: contracts.ProviderS3Compatible,
			wantRule:     RuleS3CompatibleFallback,
			wantHost:     "custom-storage.internal.net",
			wantBucket:   "",
		},
		{
			name:         "schemeless aws with query containing scheme separator",
			endpoint:     "user:secretpass@s3.us-west-2.amazonaws.com/x?r=a://b",
			bucketURL:    "",
			wantProvider: contracts.ProviderAWS,
			wantRule:     RuleAWSHostname,
			wantHost:     "s3.us-west-2.amazonaws.com",
			wantBucket:   "",
		},
		{
			name:         "schemeless gcs with query containing scheme separator",
			endpoint:     "",
			bucketURL:    "user:secretpass@sample-bucket.storage.googleapis.com/x?r=a://b",
			wantProvider: contracts.ProviderGCS,
			wantRule:     RuleGCSHostname,
			wantHost:     "sample-bucket.storage.googleapis.com",
			wantBucket:   "sample-bucket",
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

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
		name      string
		endpoint  string
		bucketURL string
		expected  contracts.Provider
	}{
		{
			name:      "endpoint GCS overrides bucketURL AWS",
			endpoint:  "https://storage.googleapis.com",
			bucketURL: "https://s3.amazonaws.com",
			expected:  contracts.ProviderGCS,
		},
		{
			name:      "endpoint AWS overrides bucketURL GCS",
			endpoint:  "https://s3.amazonaws.com",
			bucketURL: "https://storage.googleapis.com",
			expected:  contracts.ProviderAWS,
		},
		{
			name:      "endpoint custom overrides bucketURL AWS",
			endpoint:  "http://localhost:9000",
			bucketURL: "https://s3.amazonaws.com",
			expected:  contracts.ProviderS3Compatible,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := DetectProvider(tc.endpoint, tc.bucketURL)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, p)
		})
	}
}

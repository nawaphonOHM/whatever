package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestProvider_StringAndValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		provider contracts.Provider
		expected string
		valid    bool
	}{
		{
			name:     "AWS provider",
			provider: contracts.ProviderAWS,
			expected: "aws",
			valid:    true,
		},
		{
			name:     "GCS provider",
			provider: contracts.ProviderGCS,
			expected: "gcs",
			valid:    true,
		},
		{
			name:     "S3-compatible provider",
			provider: contracts.ProviderS3Compatible,
			expected: "s3-compatible",
			valid:    true,
		},
		{
			name:     "Empty provider",
			provider: contracts.Provider(""),
			expected: "",
			valid:    false,
		},
		{
			name:     "Unsupported provider",
			provider: contracts.Provider("azure-blob"),
			expected: "azure-blob",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.provider.String())
			assert.Equal(t, tt.valid, tt.provider.IsValid())
		})
	}
}

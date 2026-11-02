package client_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func assertErrorRedacts(t *testing.T, err error, forbidden []string) {
	t.Helper()
	if !assert.Error(t, err) {
		return
	}
	for _, secret := range forbidden {
		assert.NotContains(t, err.Error(), secret)
	}
}

func TestBuildS3Options_CredentialRedactionOnError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cfg       *config.Config
		provider  contracts.Provider
		name      string
		forbidden []string
	}{
		{
			name: "AWS-style secret with slashes in endpoint",
			cfg: &config.Config{
				Bucket:   testBucket,
				Endpoint: "https://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@minio.local:9000/x",
			},
			provider:  contracts.ProviderS3Compatible,
			forbidden: []string{"AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI", "K7MDENG", "bPxRfiCYEXAMPLEKEY"},
		},
		{
			name: "Nested scheme credentials in endpoint path",
			cfg: &config.Config{
				Bucket:   testBucket,
				Endpoint: "https://user:8080/supersecret/https://key@s3.amazonaws.com",
			},
			provider:  contracts.ProviderAWS,
			forbidden: []string{"supersecret", "key@s3.amazonaws.com"},
		},
		{
			name: "Credential-bearing BucketURL path mismatch",
			cfg: &config.Config{
				Bucket:    testBucket,
				BucketURL: "https://AKIAKEY:supersecret@storage.googleapis.com/otherbucket",
			},
			provider:  contracts.ProviderGCS,
			forbidden: []string{"supersecret", "AKIAKEY"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, err := client.BuildS3Options(tt.cfg, tt.provider, nil)
			assert.Equal(t, s3.Options{}, opts)
			assertErrorRedacts(t, err, tt.forbidden)
		})
	}
}

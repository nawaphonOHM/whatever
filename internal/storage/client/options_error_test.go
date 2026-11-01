package client_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_RegionConflictError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Region:   "us-west-2",
		Endpoint: "https://s3.ap-southeast-1.amazonaws.com",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "region mismatch")
	assert.Equal(t, s3.Options{}, opts)
}

func TestBuildS3Options_BucketURLSubdomainMismatchError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://otherbucket.s3.amazonaws.com",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
	assert.Equal(t, s3.Options{}, opts)
}

func TestBuildS3Options_BucketURLPathMismatchError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://storage.googleapis.com/otherbucket",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not match")
	assert.Equal(t, s3.Options{}, opts)
}

func TestBuildS3Options_InvalidTargetURLError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "http://invalid host with space",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	assert.Error(t, err)
	assert.Equal(t, s3.Options{}, opts)
}

func TestBuildS3Options_InvalidPortError(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Endpoint: "http://localhost:99999",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
	assert.Error(t, err)
	assert.Equal(t, s3.Options{}, opts)
}

func TestBuildS3Options_CredentialRedactionOnError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		cfg       *config.Config
		provider  contracts.Provider
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
			name: "Unencoded userinfo separator with port in endpoint",
			cfg: &config.Config{
				Bucket:   testBucket,
				Endpoint: "https://user:8080/supersecret/key@s3.amazonaws.com",
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
			assert.Error(t, err)
			assert.Equal(t, s3.Options{}, opts)
			for _, secret := range tt.forbidden {
				assert.NotContains(t, err.Error(), secret)
			}
		})
	}
}

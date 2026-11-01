package client_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_GCSEndpoints_Checksum(t *testing.T) {
	t.Parallel()

	endpoints := []string{
		"https://storage.googleapis.com",
		"https://mybucket.storage.googleapis.com",
		"https://storage.us-central1.rep.googleapis.com",
		"https://storage-psc.p.googleapis.com",
	}

	for _, ep := range endpoints {
		cfg := &config.Config{Bucket: testBucket, Endpoint: ep}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
		require.NoError(t, err)
		assertWhenRequiredChecksum(t, opts)
	}
}

func TestBuildS3Options_GCSBucketURLs_Checksum(t *testing.T) {
	t.Parallel()

	bucketURLs := []string{
		"https://storage.googleapis.com/mybucket",
		"https://mybucket.storage.googleapis.com",
		"https://storage.us-central1.rep.googleapis.com/mybucket",
		"https://mybucket.storage.us-central1.rep.googleapis.com",
		"https://storage-psc.p.googleapis.com/mybucket",
		"https://mybucket.storage-psc.p.googleapis.com",
	}

	for _, bURL := range bucketURLs {
		cfg := &config.Config{BucketURL: bURL}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
		require.NoError(t, err)
		assertWhenRequiredChecksum(t, opts)
	}
}

func TestBuildS3Options_Checksum_DriverOverride(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{Bucket: testBucket}
	optsObj := client.NewOptions(
		client.WithDriverOptions(func(o *s3.Options) {
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenSupported
		}),
	)

	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, optsObj)
	require.NoError(t, err)
	assert.Equal(t, aws.RequestChecksumCalculationWhenSupported, opts.RequestChecksumCalculation)
	assert.Equal(t, aws.ResponseChecksumValidationWhenRequired, opts.ResponseChecksumValidation)
	assert.True(t, opts.DisableLogOutputChecksumValidationSkipped)
}

package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const testAWSRegionEU = "eu-west-1"

func TestBuildS3Options_AWSAmazonEndpointTrailingDot(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cfg        *config.Config
		name       string
		wantRegion string
	}{
		{
			name: "global endpoint with configured region",
			cfg: &config.Config{
				Bucket:   testBucket,
				Region:   testAWSRegionSE,
				Endpoint: "https://s3.amazonaws.com.",
			},
			wantRegion: testAWSRegionSE,
		},
		{
			name: "regional endpoint without configured region",
			cfg: &config.Config{
				Bucket:   testBucket,
				Endpoint: "https://s3.eu-west-1.amazonaws.com.",
			},
			wantRegion: testAWSRegionEU,
		},
		{
			name: "regional bucketURL without configured region",
			cfg: &config.Config{
				Bucket:    testBucket,
				BucketURL: "https://mybucket.s3.eu-west-1.amazonaws.com.",
			},
			wantRegion: testAWSRegionEU,
		},
		{
			name: "regional endpoint with matching configured region",
			cfg: &config.Config{
				Bucket:   testBucket,
				Region:   testAWSRegionEU,
				Endpoint: "https://s3.eu-west-1.amazonaws.com.",
			},
			wantRegion: testAWSRegionEU,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts, err := client.BuildS3Options(tt.cfg, contracts.ProviderAWS, nil)
			require.NoError(t, err)
			assert.Nil(t, opts.BaseEndpoint)
			assert.False(t, opts.UsePathStyle)
			assert.Equal(t, tt.wantRegion, opts.Region)
		})
	}
}

func TestBuildS3Options_UserinfoStrippedEndpoint(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Endpoint: "https://user:pass@minio.internal:9000",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, "https://minio.internal:9000", *opts.BaseEndpoint)
	assert.NotContains(t, *opts.BaseEndpoint, "user")
	assert.NotContains(t, *opts.BaseEndpoint, "pass")
}

func TestBuildS3Options_UserinfoStrippedBucketURL(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://admin:secret@storage.googleapis.com/" + testBucket,
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, gcsServiceEndpoint, *opts.BaseEndpoint)
}

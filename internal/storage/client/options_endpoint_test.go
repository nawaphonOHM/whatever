package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const gcsServiceEndpoint = "https://storage.googleapis.com"

func TestBuildS3Options_AWSAmazonEndpointSkipped(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Region:   testAWSRegionSE,
		Endpoint: "https://s3.amazonaws.com",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	require.NoError(t, err)
	assert.Nil(t, opts.BaseEndpoint)
	assert.False(t, opts.UsePathStyle)
	assert.Equal(t, testAWSRegionSE, opts.Region)
}

func TestBuildS3Options_CustomS3LabelKeepsEndpoint(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Endpoint: "https://s3.minio.internal.net",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, "https://s3.minio.internal.net", *opts.BaseEndpoint)
}

func TestBuildS3Options_BucketURLAWSVirtual(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://mybucket.s3.us-east-1.amazonaws.com",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	require.NoError(t, err)
	assert.Nil(t, opts.BaseEndpoint)
}

func TestBuildS3Options_BucketURLGCSPath(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: gcsServiceEndpoint + "/" + testBucket,
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, gcsServiceEndpoint, *opts.BaseEndpoint)
	assert.True(t, opts.UsePathStyle)
}

func TestBuildS3Options_BucketURLGCSVirtual(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://mybucket.storage.googleapis.com",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, gcsServiceEndpoint, *opts.BaseEndpoint)
}

func TestBuildS3Options_BucketURLMinioPath(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "http://localhost:9000/mybucket",
		Insecure:  true,
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, "http://localhost:9000", *opts.BaseEndpoint)
	assert.True(t, opts.UsePathStyle)
}

func TestBuildS3Options_AWSAmazonEndpointTrailingDot(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Region:   testAWSRegionSE,
		Endpoint: "https://s3.amazonaws.com.",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
	require.NoError(t, err)
	assert.Nil(t, opts.BaseEndpoint)
	assert.False(t, opts.UsePathStyle)
}

func TestBuildS3Options_UserinfoStrippedEndpoint(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:   testBucket,
		Endpoint: "https://user:password@minio.internal:9000",
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, "https://minio.internal:9000", *opts.BaseEndpoint)
	assert.NotContains(t, *opts.BaseEndpoint, "user")
	assert.NotContains(t, *opts.BaseEndpoint, "password")
	assert.True(t, opts.UsePathStyle)
}

func TestBuildS3Options_UserinfoStrippedBucketURL(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{
		Bucket:    testBucket,
		BucketURL: "https://user:password@storage.googleapis.com/" + testBucket,
	}
	opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
	require.NoError(t, err)
	require.NotNil(t, opts.BaseEndpoint)
	assert.Equal(t, gcsServiceEndpoint, *opts.BaseEndpoint)
	assert.NotContains(t, *opts.BaseEndpoint, "user")
	assert.NotContains(t, *opts.BaseEndpoint, "password")
	assert.True(t, opts.UsePathStyle)
}

func TestBuildS3Options_InsecureAndPortFormatting(t *testing.T) {
	t.Parallel()

	t.Run("insecure applies http scheme", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   testBucket,
			Endpoint: "minio.internal.net",
			Insecure: true,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "http://minio.internal.net", *opts.BaseEndpoint)
	})

	t.Run("disable ssl applies http scheme", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:     testBucket,
			Endpoint:   "minio.internal.net",
			DisableSSL: true,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "http://minio.internal.net", *opts.BaseEndpoint)
	})

	t.Run("port appended when missing", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   testBucket,
			Endpoint: "https://custom.s3.internal",
			Port:     8443,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "https://custom.s3.internal:8443", *opts.BaseEndpoint)
	})

	t.Run("explicit use path style on aws", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:       testBucket,
			Endpoint:     "https://s3.amazonaws.com",
			UsePathStyle: true,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Nil(t, opts.BaseEndpoint)
		assert.True(t, opts.UsePathStyle)
	})

	t.Run("aws insecure keeps explicit http base endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   testBucket,
			Endpoint: "https://s3.amazonaws.com",
			Insecure: true,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "http://s3.amazonaws.com", *opts.BaseEndpoint)
	})

	t.Run("aws custom port keeps explicit base endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   testBucket,
			Endpoint: "https://s3.amazonaws.com",
			Port:     8443,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "https://s3.amazonaws.com:8443", *opts.BaseEndpoint)
	})

	t.Run("aws path prefix keeps explicit base endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   testBucket,
			Endpoint: "https://s3.amazonaws.com/custom/prefix",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "https://s3.amazonaws.com/custom/prefix", *opts.BaseEndpoint)
	})
}

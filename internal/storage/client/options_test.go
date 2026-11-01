package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestBuildS3Options_Providers(t *testing.T) {
	t.Parallel()

	t.Run("aws default endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket: "test-b",
			Region: "us-west-2",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Equal(t, "us-west-2", opts.Region)
		assert.Nil(t, opts.BaseEndpoint)
		assert.False(t, opts.UsePathStyle)
		assertWhenSupportedChecksum(t, opts)
	})

	t.Run("gcs default endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{Bucket: "test-b"}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "https://storage.googleapis.com", *opts.BaseEndpoint)
		assert.True(t, opts.UsePathStyle)
		assertWhenRequiredChecksum(t, opts)
	})

	t.Run("s3 compatible custom endpoint and port", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Bucket:   "test-b",
			Endpoint: "localhost",
			Port:     testMinioPort,
			Insecure: true,
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "http://localhost:9000", *opts.BaseEndpoint)
		assert.True(t, opts.UsePathStyle)
		assertWhenRequiredChecksum(t, opts)
	})
}

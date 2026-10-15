package client_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	})

	t.Run("gcs default endpoint", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{Bucket: "test-b"}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
		require.NoError(t, err)
		require.NotNil(t, opts.BaseEndpoint)
		assert.Equal(t, "https://storage.googleapis.com", *opts.BaseEndpoint)
		assert.True(t, opts.UsePathStyle)
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
	})
}

const (
	testHTTPTimeout    = 5 * time.Second
	testRequestTimeout = 10 * time.Second
)

func assertCustomCredentials(t *testing.T, opts s3.Options, customHTTP *http.Client) {
	t.Helper()
	assert.Same(t, customHTTP, opts.HTTPClient)
	assert.True(t, opts.UsePathStyle)
	require.NotNil(t, opts.Credentials)
	creds, err := opts.Credentials.Retrieve(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "ak", creds.AccessKeyID)
}

func TestBuildS3Options_CredentialsAndHTTP(t *testing.T) {
	t.Parallel()

	customHTTP := &http.Client{Timeout: testHTTPTimeout}
	cfg := &config.Config{
		Bucket:          "b",
		AccessKeyID:     "ak",
		SecretAccessKey: "sk",
		SessionToken:    "tok",
		RequestTimeout:  testRequestTimeout,
	}
	optsObj := client.NewOptions(
		client.WithHTTPClient(customHTTP),
		client.WithSkipProbe(true),
		client.WithDriverOptions(func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	)

	opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, optsObj)
	require.NoError(t, err)
	assertCustomCredentials(t, opts, customHTTP)
}

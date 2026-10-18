package storage_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const testPresignOptionExpiry time.Duration = 5

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func isolatePublicAWSEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AWS_CONFIG_FILE", "/nonexistent/path/config")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/nonexistent/path/credentials")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_ENDPOINT_URL", "")
	t.Setenv("AWS_ENDPOINT_URL_S3", "")
}

func setPublicConnectEnv(t *testing.T) {
	t.Helper()
	isolatePublicAWSEnv(t)
	t.Setenv("OHM9996_STORAGE_BUCKET", "public-bucket")
	t.Setenv("OHM9996_STORAGE_REGION", "us-east-1")
	t.Setenv("OHM9996_STORAGE_ENDPOINT", "https://s3.amazonaws.com")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
}

func publicConnectOpts(client *http.Client) []storage.Option {
	return []storage.Option{
		storage.WithHTTPClient(client),
		storage.WithSkipProbe(false),
		storage.WithDriverOptions(func(opts *s3.Options) { opts.UsePathStyle = true }),
		storage.WithPresignOptions(func(opts *s3.PresignOptions) { opts.Expires = testPresignOptionExpiry }),
	}
}

func TestPublicStorage_Connect(t *testing.T) {
	setPublicConnectEnv(t)
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Header:     make(http.Header),
			}, nil
		}),
	}
	c, err := storage.Connect(context.Background(), publicConnectOpts(client)...)
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, "public-bucket", c.DefaultBucket())
	assert.Equal(t, storage.ProviderAWS, c.Provider())
}

func customGCSConfig() *storage.Config {
	cfg := storage.DefaultConfig()
	cfg.Bucket = "custom-config-bucket"
	cfg.Endpoint = "https://storage.googleapis.com"
	cfg.AccessKeyID = "AKIAIOSFODNN7EXAMPLE"
	cfg.SecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	return cfg
}

func TestPublicStorage_ConnectWithConfig(t *testing.T) {
	t.Parallel()

	c, err := storage.ConnectWithConfig(
		context.Background(),
		customGCSConfig(),
		storage.WithSkipProbe(true),
	)

	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, "custom-config-bucket", c.DefaultBucket())
	assert.Equal(t, storage.ProviderGCS, c.Provider())
}

func TestPublicStorage_ImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ storage.Storage = (*storage.Client)(nil)
}

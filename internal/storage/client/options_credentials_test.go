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

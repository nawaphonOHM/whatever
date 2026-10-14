package client_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

const (
	testHTTPEndpoint  = "http://localhost:9000"
	testHTTPSEndpoint = "https://minio.internal.net"
)

func disableRetries(o *s3.Options) {
	o.RetryMaxAttempts = 1
}

func connectOptions(httpClient *http.Client) []client.Option {
	return []client.Option{
		client.WithHTTPClient(httpClient),
		client.WithSkipProbe(true),
		client.WithDriverOptions(disableRetries),
	}
}

func connectOptionsFn(fn func(*http.Request) (*http.Response, error)) []client.Option {
	return connectOptions(newMockHTTPClient(fn))
}

func endpointConfig(endpoint string) *config.Config {
	cfg := testConfigWithKeys()
	cfg.Endpoint = endpoint
	return cfg
}

func connectRecorded(t *testing.T, cfg *config.Config, rec *recorder) *client.Client {
	t.Helper()
	c, err := client.ConnectWithConfig(t.Context(), cfg, connectOptions(rec.httpClient())...)
	require.NoError(t, err)
	return c
}

package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestClient_Ping_Success(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodHead, req.Method)
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	err := c.Ping(context.Background())
	require.NoError(t, err)
}

func TestClient_Ping_BucketNotFound(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return newHTTPResponse(http.StatusNotFound, "", nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	err := c.Ping(context.Background())
	assert.ErrorIs(t, err, contracts.ErrBucketNotFound)
}

func TestClient_Ping_NilClientAndEmptyBucket(t *testing.T) {
	t.Parallel()

	var c *client.Client
	assert.ErrorIs(t, c.Ping(context.Background()), contracts.ErrNilClient)

	emptyBucketClient := client.NewClient(client.ClientParams{
		Raw:           s3.New(s3.Options{Region: "us-east-1"}),
		DefaultBucket: "",
	})
	assert.ErrorIs(t, emptyBucketClient.Ping(context.Background()), contracts.ErrBucketNotFound)
}

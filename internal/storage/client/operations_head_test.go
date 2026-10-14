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

const expectedContentLength int64 = 1024

func TestClient_HeadObject_Success(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodHead, req.Method)
		return newHTTPResponse(http.StatusOK, "", map[string]string{
			"Content-Type":   "image/png",
			"Content-Length": "1024",
			"ETag":           `"head-etag"`,
		}), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.HeadObject(context.Background(), contracts.HeadObjectInput{
		Key: "image.png",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, expectedContentLength, out.ContentLength)
	assert.Equal(t, "image/png", out.ContentType)
	assert.Equal(t, `"head-etag"`, out.ETag)
}

func TestClient_Exists(t *testing.T) {
	t.Parallel()

	t.Run("exists true on 200", func(t *testing.T) {
		t.Parallel()
		httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			_ = req
			return newHTTPResponse(http.StatusOK, "", nil), nil
		})
		c := client.NewClient(client.ClientParams{
			Raw:           s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient}),
			DefaultBucket: "my-bucket",
		})
		exists, err := c.Exists(context.Background(), "", "existing.txt")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("exists false on 404", func(t *testing.T) {
		t.Parallel()
		httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/" {
				return newHTTPResponse(http.StatusOK, "", nil), nil
			}
			return newHTTPResponse(http.StatusNotFound, "", nil), nil
		})
		c := client.NewClient(client.ClientParams{
			Raw:           s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient}),
			DefaultBucket: "my-bucket",
		})
		exists, err := c.Exists(context.Background(), "", "not-found.txt")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("nil client returns error", func(t *testing.T) {
		t.Parallel()
		var c *client.Client
		exists, err := c.Exists(context.Background(), "b", "k")
		assert.False(t, exists)
		assert.ErrorIs(t, err, contracts.ErrNilClient)
	})
}

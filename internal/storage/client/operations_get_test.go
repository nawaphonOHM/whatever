package client_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const noSuchKeyXML = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>`

func assertGetObjectOutput(t *testing.T, out *contracts.GetObjectOutput) {
	t.Helper()
	defer func() {
		assert.NoError(t, out.Body.Close())
	}()

	body, err := io.ReadAll(out.Body)
	require.NoError(t, err)
	assert.Equal(t, "file content payload", string(body))
	assert.Equal(t, "application/json", out.ContentType)
	assert.Equal(t, `"etag-456"`, out.ETag)
	assert.Equal(t, "1.0", out.Metadata["version"])
}

func TestClient_GetObject_Success(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		return newHTTPResponse(http.StatusOK, "file content payload", map[string]string{
			"Content-Type":       "application/json",
			"Content-Length":     "20",
			"ETag":               `"etag-456"`,
			"Last-Modified":      "Wed, 21 Oct 2025 07:28:00 GMT",
			"x-amz-meta-version": "1.0",
		}), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.GetObject(context.Background(), contracts.GetObjectInput{
		Key: "file.json",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assertGetObjectOutput(t, out)
}

func TestClient_GetObject_NotFound(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return newHTTPResponse(http.StatusNotFound, noSuchKeyXML, nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.GetObject(context.Background(), contracts.GetObjectInput{
		Key: "non-existent.txt",
	})

	assert.Nil(t, out)
	assert.ErrorIs(t, err, contracts.ErrObjectNotFound)
}

func TestClient_GetObject_NilClient(t *testing.T) {
	t.Parallel()

	var c *client.Client
	out, err := c.GetObject(context.Background(), contracts.GetObjectInput{
		Key: "file.txt",
	})
	assert.Nil(t, out)
	assert.ErrorIs(t, err, contracts.ErrNilClient)
}

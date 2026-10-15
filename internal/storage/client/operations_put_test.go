package client_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestClient_PutObject_Success(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPut, req.Method)
		assert.Contains(t, req.URL.Path, "test-key.txt")
		assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
		assert.Equal(t, "custom-value", req.Header.Get("x-amz-meta-custom-meta"))

		return newHTTPResponse(http.StatusOK, "", map[string]string{
			"ETag":             `"test-etag-123"`,
			"x-amz-version-id": "v1.0",
		}), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.PutObject(context.Background(), contracts.PutObjectInput{
		Key:         "test-key.txt",
		Body:        bytes.NewBufferString("hello world"),
		ContentType: "text/plain",
		Metadata:    map[string]string{"custom-meta": "custom-value"},
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, `"test-etag-123"`, out.ETag)
	assert.Equal(t, "v1.0", out.VersionID)
}

func TestClient_PutObject_NilClient(t *testing.T) {
	t.Parallel()

	var c *client.Client
	out, err := c.PutObject(context.Background(), contracts.PutObjectInput{
		Key: "file.txt",
	})
	assert.Nil(t, out)
	assert.ErrorIs(t, err, contracts.ErrNilClient)
}

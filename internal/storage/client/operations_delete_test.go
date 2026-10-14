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

func TestClient_DeleteObject_Success(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodDelete, req.Method)
		assert.Contains(t, req.URL.Path, "to-delete.txt")
		return newHTTPResponse(http.StatusNoContent, "", map[string]string{
			"x-amz-version-id": "del-v2",
		}), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.DeleteObject(context.Background(), contracts.DeleteObjectInput{
		Key: "to-delete.txt",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "del-v2", out.VersionID)
}

func TestClient_DeleteObject_BucketOverride(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodDelete, req.Method)
		assert.Contains(t, req.URL.Path, "override-bucket")
		return newHTTPResponse(http.StatusNoContent, "", nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient, UsePathStyle: true})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "default-bucket",
	})

	out, err := c.DeleteObject(context.Background(), contracts.DeleteObjectInput{
		Bucket: "override-bucket",
		Key:    "file.txt",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
}

func TestClient_DeleteObject_NilClient(t *testing.T) {
	t.Parallel()

	var c *client.Client
	out, err := c.DeleteObject(context.Background(), contracts.DeleteObjectInput{
		Key: "to-delete.txt",
	})
	assert.Nil(t, out)
	assert.ErrorIs(t, err, contracts.ErrNilClient)
}

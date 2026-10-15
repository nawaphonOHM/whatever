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

const (
	expectedMultiPageObjectCount = 3
	expectedFirstObjectSize      = int64(100)
	expectedRequestCount         = 2
)

func assertMultiPageListOutput(t *testing.T, out *contracts.ListObjectsOutput, requestCount int) {
	t.Helper()
	assert.Equal(t, expectedRequestCount, requestCount)
	assert.Len(t, out.Objects, expectedMultiPageObjectCount)
	assert.Equal(t, "photos/1.jpg", out.Objects[0].Key)
	assert.Equal(t, expectedFirstObjectSize, out.Objects[0].Size)
	assert.Equal(t, "photos/3.jpg", out.Objects[2].Key)
	assert.Equal(t, []string{"photos/2025/"}, out.Prefixes)
}

func TestClient_ListObjects_MultiPage(t *testing.T) {
	t.Parallel()

	requestCount := 0
	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		requestCount++
		if req.URL.Query().Get("continuation-token") == "token-2" {
			return newHTTPResponse(http.StatusOK, listXMLPage2, nil), nil
		}
		return newHTTPResponse(http.StatusOK, listXMLPage1, nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.ListObjects(context.Background(), contracts.ListObjectsInput{
		Prefix:    "photos/",
		Delimiter: "/",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assertMultiPageListOutput(t, out, requestCount)
}

func TestClient_ListObjects_MaxKeysEarlyExit(t *testing.T) {
	t.Parallel()

	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		_ = req
		return newHTTPResponse(http.StatusOK, listXMLPage1, nil), nil
	})

	raw := s3.New(s3.Options{Region: "us-east-1", RetryMaxAttempts: 1, HTTPClient: httpClient})
	c := client.NewClient(client.ClientParams{
		Raw:           raw,
		Provider:      contracts.ProviderAWS,
		DefaultBucket: "my-bucket",
	})

	out, err := c.ListObjects(context.Background(), contracts.ListObjectsInput{
		Prefix:  "photos/",
		MaxKeys: 1,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Len(t, out.Objects, 1)
	assert.Equal(t, "photos/1.jpg", out.Objects[0].Key)
}

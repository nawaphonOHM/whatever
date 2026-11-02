package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func newTestGCSClient(t *testing.T) *client.Client {
	t.Helper()
	cfg := testConfigWithKeys()
	cfg.Endpoint = gcsServiceEndpoint
	c, err := client.ConnectWithConfig(
		context.Background(), cfg, client.WithSkipProbe(true),
	)
	require.NoError(t, err)
	return c
}

func assertPresignedURL(t *testing.T, out *contracts.PresignedURL, method, key string) {
	t.Helper()
	require.NotNil(t, out)
	assert.Equal(t, method, out.Method)
	assert.Contains(t, out.URL, gcsServiceEndpoint)
	assert.Contains(t, out.URL, key)
	assert.Contains(t, out.URL, "X-Amz-Signature")
}

func TestConnect_GCS_PresignGetObject(t *testing.T) {
	t.Parallel()
	c := newTestGCSClient(t)
	out, err := c.PresignGetObject(context.Background(), contracts.PresignGetObjectInput{
		Key:       "download.pdf",
		ExpiresIn: testPresignExpiry5m,
	})
	require.NoError(t, err)
	assertPresignedURL(t, out, http.MethodGet, "download.pdf")
	assert.Equal(t, testPresignExpiry5m, out.ExpiresIn)
}

func TestConnect_GCS_PresignPutObject(t *testing.T) {
	t.Parallel()
	c := newTestGCSClient(t)
	out, err := c.PresignPutObject(context.Background(), contracts.PresignPutObjectInput{
		Key:         "upload.png",
		ContentType: "image/png",
		ExpiresIn:   testPresignExpiry20m,
	})
	require.NoError(t, err)
	assertPresignedURL(t, out, http.MethodPut, "upload.png")
	assert.Equal(t, []string{"image/png"}, out.Headers["Content-Type"])
	assert.Equal(t, testPresignExpiry20m, out.ExpiresIn)
}

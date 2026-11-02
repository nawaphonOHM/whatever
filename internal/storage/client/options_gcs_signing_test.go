package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type signingTestCapture struct {
	auth           string
	acceptEncoding string
}

func extractSignedHeaders(auth string) string {
	for part := range strings.SplitSeq(auth, ",") {
		if after, ok := strings.CutPrefix(strings.TrimSpace(part), "SignedHeaders="); ok {
			return after
		}
	}
	return ""
}

func headWithCapture(t *testing.T, endpoint string) signingTestCapture {
	t.Helper()
	var rec signingTestCapture
	httpClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		rec.auth = req.Header.Get("Authorization")
		rec.acceptEncoding = req.Header.Get("Accept-Encoding")
		return newHTTPResponse(http.StatusOK, "", nil), nil
	})
	cfg := testConfigWithKeys()
	cfg.Endpoint = endpoint
	c, err := client.ConnectWithConfig(
		context.Background(), cfg, client.WithHTTPClient(httpClient), client.WithSkipProbe(true),
	)
	require.NoError(t, err)
	_, err = c.HeadObject(context.Background(), contracts.HeadObjectInput{Key: testObjectKey})
	require.NoError(t, err)
	return rec
}

func TestConnect_GCS_ExcludesAcceptEncodingFromSignedHeaders(t *testing.T) {
	t.Parallel()
	got := headWithCapture(t, gcsServiceEndpoint)
	require.NotEmpty(t, got.auth)
	assert.Contains(t, got.auth, sigV4Algorithm)
	signedHeaders := extractSignedHeaders(got.auth)
	require.NotEmpty(t, signedHeaders)
	assert.NotContains(t, strings.Split(signedHeaders, ";"), "accept-encoding")
	assert.Equal(t, "identity", got.acceptEncoding)
}

func TestConnect_AWS_IncludesSignedHeadersNormally(t *testing.T) {
	t.Parallel()
	got := headWithCapture(t, "https://s3.us-east-1.amazonaws.com")
	require.NotEmpty(t, got.auth)
	assert.Contains(t, got.auth, sigV4Algorithm)
	signedHeaders := extractSignedHeaders(got.auth)
	require.NotEmpty(t, signedHeaders)
	assert.Contains(t, strings.Split(signedHeaders, ";"), "accept-encoding")
	assert.Equal(t, "identity", got.acceptEncoding)
}

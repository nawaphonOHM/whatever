//go:build testcontainers

package storage_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const presignedHTTPTimeout = 15 * time.Second

func sendLocalStackPresignedRequest(t *testing.T, signed *storage.PresignedURL, payload []byte) (int, []byte) {
	t.Helper()
	request := newLocalStackPresignedRequest(t, signed, payload)
	response, err := (&http.Client{Timeout: presignedHTTPTimeout}).Do(request)
	require.NoError(t, err)
	defer func() { require.NoError(t, response.Body.Close()) }()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response.StatusCode, body
}

func newLocalStackPresignedRequest(t *testing.T, signed *storage.PresignedURL, payload []byte) *http.Request {
	request, err := http.NewRequestWithContext(t.Context(), signed.Method, signed.URL, bytes.NewReader(payload))
	require.NoError(t, err)
	copyLocalStackPresignedHeaders(request, signed.Headers)
	return request
}

func copyLocalStackPresignedHeaders(request *http.Request, headers http.Header) {
	for name, values := range headers {
		if strings.EqualFold(name, "host") {
			continue
		}
		request.Header[name] = append([]string(nil), values...)
	}
}

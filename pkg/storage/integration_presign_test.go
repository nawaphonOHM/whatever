//go:build testcontainers

package storage_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const presignedURLExpiry = 5 * time.Minute

func TestLocalStackPresignedHTTPRoundTrip(t *testing.T) {
	requireSharedLocalStack(t)
	key := localStackObjectKey(t, "presigned")
	cleanupLocalStackObject(t, key)
	payload := []byte("uploaded through a presigned URL")
	put := createLocalStackPresignedPut(t, key)
	uploadLocalStackPresignedPayload(t, put, payload)
	get := createLocalStackPresignedGet(t, key)
	assertLocalStackPresignedPayload(t, get, payload)
}

func createLocalStackPresignedPut(t *testing.T, key string) *storage.PresignedURL {
	put, err := sharedStorage.PresignPutObject(t.Context(), storage.PresignPutObjectInput{
		Key: key, ContentType: "text/plain", ExpiresIn: presignedURLExpiry,
	})
	require.NoError(t, err)
	require.Equal(t, http.MethodPut, put.Method)
	require.Equal(t, []string{"text/plain"}, put.Headers["Content-Type"])
	return put
}

func createLocalStackPresignedGet(t *testing.T, key string) *storage.PresignedURL {
	get, err := sharedStorage.PresignGetObject(t.Context(), storage.PresignGetObjectInput{
		Key: key, ExpiresIn: presignedURLExpiry,
	})
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, get.Method)
	return get
}

func uploadLocalStackPresignedPayload(t *testing.T, signed *storage.PresignedURL, payload []byte) {
	status, _ := sendLocalStackPresignedRequest(t, signed, payload)
	require.Equal(t, http.StatusOK, status)
}

func assertLocalStackPresignedPayload(t *testing.T, signed *storage.PresignedURL, want []byte) {
	status, got := sendLocalStackPresignedRequest(t, signed, nil)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, want, got)
}

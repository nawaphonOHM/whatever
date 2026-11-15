//go:build testcontainers

package storage_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const multipartStreamSize = 8*1024*1024 + 1

func TestLocalStackStreamUploads(t *testing.T) {
	requireSharedLocalStack(t)
	t.Run("seekable", func(t *testing.T) {
		payload := []byte("seekable stream payload")
		uploadAndVerifyLocalStackObject(t, strings.NewReader(string(payload)), payload)
	})
	t.Run("non-seekable", func(t *testing.T) {
		payload := []byte("non-seekable stream payload")
		uploadAndVerifyLocalStackObject(t, bytes.NewBuffer(payload), payload)
	})
	t.Run("multipart-sized non-seekable", func(t *testing.T) {
		payload := bytes.Repeat([]byte("m"), multipartStreamSize)
		uploadAndVerifyLocalStackObject(t, bytes.NewBuffer(payload), payload)
	})
}

func uploadAndVerifyLocalStackObject(t *testing.T, body io.Reader, want []byte) {
	t.Helper()
	key := localStackObjectKey(t, "stream")
	_, err := sharedStorage.PutObject(context.Background(), storage.PutObjectInput{Key: key, Body: body})
	require.NoError(t, err)
	cleanupLocalStackObject(t, key)
	object, err := sharedStorage.GetObject(context.Background(), storage.GetObjectInput{Key: key})
	require.NoError(t, err)
	require.Equal(t, want, readAndCloseLocalStackObjectBody(t, object.Body))
}

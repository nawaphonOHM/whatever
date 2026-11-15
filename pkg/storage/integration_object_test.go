//go:build testcontainers

package storage_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func localStackObjectKey(t *testing.T, suffix string) string {
	t.Helper()
	return fmt.Sprintf("integration/%s/%s", strings.ReplaceAll(t.Name(), "/", "-"), suffix)
}

func cleanupLocalStackObject(t *testing.T, key string) {
	t.Helper()
	t.Cleanup(func() {
		if _, err := sharedStorage.DeleteObject(context.Background(), storage.DeleteObjectInput{Key: key}); err != nil {
			t.Errorf("delete integration object %q during cleanup: %v", key, err)
		}
	})
}

func readAndCloseLocalStackObjectBody(t *testing.T, body io.ReadCloser) []byte {
	t.Helper()
	got, readErr := io.ReadAll(body)
	closeErr := body.Close()
	if readErr != nil {
		require.NoError(t, readErr)
	}
	require.NoError(t, closeErr)
	return got
}

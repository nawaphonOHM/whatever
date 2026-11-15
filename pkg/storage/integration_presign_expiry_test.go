//go:build testcontainers

package storage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const invalidPresignedURLExpiry = 7*24*time.Hour + time.Second

func TestLocalStackPresignExpiry(t *testing.T) {
	requireSharedLocalStack(t)
	key := localStackObjectKey(t, "invalid-expiry")
	assertInvalidLocalStackPresignExpiry(t, key, invalidPresignedURLExpiry)
}

func assertInvalidLocalStackPresignExpiry(t *testing.T, key string, expiry time.Duration) {
	_, err := sharedStorage.PresignGetObject(t.Context(), storage.PresignGetObjectInput{
		Key: key, ExpiresIn: expiry,
	})
	require.ErrorIs(t, err, storage.ErrInvalidExpiry)
	_, err = sharedStorage.PresignPutObject(t.Context(), storage.PresignPutObjectInput{
		Key: key, ExpiresIn: expiry,
	})
	require.ErrorIs(t, err, storage.ErrInvalidExpiry)
}

func TestLocalStackDeleteAbsentObject(t *testing.T) {
	requireSharedLocalStack(t)
	key := localStackObjectKey(t, "absent-delete")
	_, err := sharedStorage.DeleteObject(t.Context(), storage.DeleteObjectInput{Key: key})
	require.NoError(t, err)
}

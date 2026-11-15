//go:build testcontainers

package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

type localStackMissingObject struct {
	getHeadErr error
	existsErr  error
	bucket     string
	key        string
}

func TestLocalStackMissingObjectErrors(t *testing.T) {
	requireSharedLocalStack(t)
	missing := localStackMissingObject{key: localStackObjectKey(t, "missing"), getHeadErr: storage.ErrObjectNotFound}
	assertLocalStackObjectMissing(t, missing)
}

func TestLocalStackMissingBucketErrors(t *testing.T) {
	requireSharedLocalStack(t)
	missing := localStackMissingObject{
		bucket: "whatever-missing-integration-bucket", key: localStackObjectKey(t, "missing-bucket"),
		getHeadErr: storage.ErrBucketNotFound, existsErr: storage.ErrBucketNotFound,
	}
	assertLocalStackObjectMissing(t, missing)
}

func assertLocalStackObjectMissing(t *testing.T, missing localStackMissingObject) {
	assertMissingLocalStackGet(t, missing)
	assertMissingLocalStackHead(t, missing)
	assertMissingLocalStackExists(t, missing)
}

func assertMissingLocalStackGet(t *testing.T, missing localStackMissingObject) {
	_, err := sharedStorage.GetObject(context.Background(), storage.GetObjectInput{
		Bucket: missing.bucket, Key: missing.key,
	})
	require.ErrorIs(t, err, missing.getHeadErr)
}

func assertMissingLocalStackHead(t *testing.T, missing localStackMissingObject) {
	_, err := sharedStorage.HeadObject(context.Background(), storage.HeadObjectInput{
		Bucket: missing.bucket, Key: missing.key,
	})
	require.ErrorIs(t, err, missing.getHeadErr)
}

func assertMissingLocalStackExists(t *testing.T, missing localStackMissingObject) {
	exists, err := sharedStorage.Exists(context.Background(), missing.bucket, missing.key)
	require.False(t, exists)
	if missing.existsErr != nil {
		require.ErrorIs(t, err, missing.existsErr)
		return
	}
	require.NoError(t, err)
}

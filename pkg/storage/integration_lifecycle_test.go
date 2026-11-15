//go:build testcontainers

package storage_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestLocalStackObjectLifecycle(t *testing.T) {
	requireSharedLocalStack(t)
	object := newLocalStackLifecycleObject(t)
	putLocalStackLifecycleObject(t, object)
	assertLocalStackObjectExists(t, object)
	assertLocalStackObjectHead(t, object)
	assertLocalStackObjectContents(t, object)
	deleteLocalStackLifecycleObject(t, object)
	assertLocalStackObjectAbsent(t, object)
}

type localStackLifecycleObject struct {
	metadata map[string]string
	key      string
	ctx      context.Context
	payload  []byte
}

func newLocalStackLifecycleObject(t *testing.T) localStackLifecycleObject {
	return localStackLifecycleObject{
		ctx: context.Background(), key: localStackObjectKey(t, "lifecycle"),
		payload:  []byte(`{"message":"LocalStack integration"}`),
		metadata: map[string]string{"source": "integration", "purpose": "lifecycle"},
	}
}

func putLocalStackLifecycleObject(t *testing.T, object localStackLifecycleObject) {
	_, err := sharedStorage.PutObject(object.ctx, storage.PutObjectInput{
		Key:         object.key,
		Body:        bytes.NewReader(object.payload),
		ContentType: "application/json",
		Metadata:    object.metadata,
	})
	require.NoError(t, err)
	cleanupLocalStackObject(t, object.key)
}

func assertLocalStackObjectExists(t *testing.T, object localStackLifecycleObject) {
	exists, err := sharedStorage.Exists(object.ctx, "", object.key)
	require.NoError(t, err)
	require.True(t, exists)
}

func assertLocalStackObjectHead(t *testing.T, object localStackLifecycleObject) {
	head, err := sharedStorage.HeadObject(object.ctx, storage.HeadObjectInput{Key: object.key})
	require.NoError(t, err)
	require.Equal(t, int64(len(object.payload)), head.ContentLength)
	require.Equal(t, "application/json", head.ContentType)
	require.Equal(t, object.metadata, head.Metadata)
	require.NotEmpty(t, head.ETag)
}

func assertLocalStackObjectContents(t *testing.T, expected localStackLifecycleObject) {
	object, err := sharedStorage.GetObject(expected.ctx, storage.GetObjectInput{Key: expected.key})
	require.NoError(t, err)
	got := readAndCloseLocalStackObjectBody(t, object.Body)
	require.Equal(t, expected.payload, got)
	require.Equal(t, int64(len(expected.payload)), object.ContentLength)
	require.Equal(t, "application/json", object.ContentType)
	require.Equal(t, expected.metadata, object.Metadata)
}

func deleteLocalStackLifecycleObject(t *testing.T, object localStackLifecycleObject) {
	_, err := sharedStorage.DeleteObject(object.ctx, storage.DeleteObjectInput{Key: object.key})
	require.NoError(t, err)
}

func assertLocalStackObjectAbsent(t *testing.T, object localStackLifecycleObject) {
	exists, err := sharedStorage.Exists(object.ctx, "", object.key)
	require.NoError(t, err)
	require.False(t, exists)
}

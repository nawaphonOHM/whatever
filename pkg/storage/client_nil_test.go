package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestPublicStorage_ConnectWithConfig_NilConfig(t *testing.T) {
	t.Parallel()

	c, err := storage.ConnectWithConfig(context.Background(), nil)
	require.Error(t, err)
	assert.Nil(t, c)
	assert.ErrorIs(t, err, storage.ErrNilConfig)
}

func TestPublicStorage_NilClient_CRUD(t *testing.T) {
	t.Parallel()

	var c *storage.Client
	ctx := context.Background()

	_, errPut := c.PutObject(ctx, storage.PutObjectInput{})
	assert.ErrorIs(t, errPut, storage.ErrNilClient)

	_, errGet := c.GetObject(ctx, storage.GetObjectInput{})
	assert.ErrorIs(t, errGet, storage.ErrNilClient)

	_, errDel := c.DeleteObject(ctx, storage.DeleteObjectInput{})
	assert.ErrorIs(t, errDel, storage.ErrNilClient)
}

func TestPublicStorage_NilClient_Inspection(t *testing.T) {
	t.Parallel()

	var c *storage.Client
	ctx := context.Background()

	_, errHead := c.HeadObject(ctx, storage.HeadObjectInput{})
	assert.ErrorIs(t, errHead, storage.ErrNilClient)

	_, errList := c.ListObjects(ctx, storage.ListObjectsInput{})
	assert.ErrorIs(t, errList, storage.ErrNilClient)

	exists, errExists := c.Exists(ctx, testBucket, testKey)
	assert.False(t, exists)
	assert.ErrorIs(t, errExists, storage.ErrNilClient)
}

func TestPublicStorage_NilClient_Presign(t *testing.T) {
	t.Parallel()

	var c *storage.Client
	ctx := context.Background()

	_, errPresignGet := c.PresignGetObject(ctx, storage.PresignGetObjectInput{})
	assert.ErrorIs(t, errPresignGet, storage.ErrNilClient)

	_, errPresignPut := c.PresignPutObject(ctx, storage.PresignPutObjectInput{})
	assert.ErrorIs(t, errPresignPut, storage.ErrNilClient)
}

func TestPublicStorage_NilClient_Accessors(t *testing.T) {
	t.Parallel()

	var c *storage.Client
	assert.ErrorIs(t, c.Ping(context.Background()), storage.ErrNilClient)
	assert.Equal(t, storage.Provider(""), c.Provider())
	assert.Equal(t, "", c.DefaultBucket())
	assert.Nil(t, c.RawClient())
	assert.Nil(t, c.PresignClient())
}

//go:build testcontainers

package azurite_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

type s3LifecycleFixture struct {
	ctx    context.Context
	client *storage.Client
	names  s3ResourceNames
}

func (f s3LifecycleFixture) verify(t *testing.T, markBucketDeleted func()) {
	t.Helper()
	want := []byte("Azurite S3 gateway integration test payload")
	f.upload(t, want)
	f.assertContents(t, want)
	f.deleteObject(t)
	f.deleteBucket(t, markBucketDeleted)
}

func (f s3LifecycleFixture) upload(t *testing.T, content []byte) {
	t.Helper()
	registerS3ObjectCleanup(t, f.client, f.names.bucket, f.names.key)
	_, err := f.client.PutObject(f.ctx, storage.PutObjectInput{
		Bucket: f.names.bucket,
		Key:    f.names.key,
		Body:   bytes.NewReader(content),
	})
	require.NoError(t, err)
}

func (f s3LifecycleFixture) assertContents(t *testing.T, want []byte) {
	t.Helper()
	download, err := f.client.GetObject(f.ctx, storage.GetObjectInput{Bucket: f.names.bucket, Key: f.names.key})
	require.NoError(t, err)
	got, readErr := io.ReadAll(download.Body)
	closeErr := download.Body.Close()
	require.NoError(t, readErr)
	require.NoError(t, closeErr)
	require.Equal(t, want, got)
}

func (f s3LifecycleFixture) deleteObject(t *testing.T) {
	t.Helper()
	_, err := f.client.DeleteObject(f.ctx, storage.DeleteObjectInput{Bucket: f.names.bucket, Key: f.names.key})
	require.NoError(t, err)
	exists, err := f.client.Exists(f.ctx, f.names.bucket, f.names.key)
	require.NoError(t, err)
	require.False(t, exists)
}

func (f s3LifecycleFixture) deleteBucket(t *testing.T, markBucketDeleted func()) {
	t.Helper()
	_, err := f.client.RawClient().DeleteBucket(f.ctx, &s3.DeleteBucketInput{Bucket: aws.String(f.names.bucket)})
	require.NoError(t, err)
	markBucketDeleted()
	_, err = f.client.RawClient().HeadBucket(f.ctx, &s3.HeadBucketInput{Bucket: aws.String(f.names.bucket)})
	require.Error(t, err)
}

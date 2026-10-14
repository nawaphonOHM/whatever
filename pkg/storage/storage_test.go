package storage_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	testBucket    = "test-bucket"
	testKey       = "test-key"
	testETag      = "\"etag-123\""
	testSize      = 42
	testLength    = 100
	testURLString = "https://example.com"
	testExpiry    = 30 * time.Minute
)

func TestPublicTypeAliasesAndConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, contracts.ProviderAWS, storage.ProviderAWS)
	assert.Equal(t, contracts.ProviderGCS, storage.ProviderGCS)
	assert.Equal(t, contracts.ProviderS3Compatible, storage.ProviderS3Compatible)

	assert.True(t, errors.Is(storage.ErrObjectNotFound, contracts.ErrObjectNotFound))
	assert.True(t, errors.Is(storage.ErrBucketNotFound, contracts.ErrBucketNotFound))
	assert.True(t, errors.Is(storage.ErrNilClient, contracts.ErrNilClient))
	assert.True(t, errors.Is(storage.ErrNilConfig, contracts.ErrNilConfig))
}

func TestPublicModels_PutAndGet(t *testing.T) {
	t.Parallel()

	putIn := storage.PutObjectInput{Bucket: testBucket, Key: testKey}
	assert.Equal(t, testBucket, putIn.Bucket)
	assert.Equal(t, testKey, putIn.Key)

	putOut := storage.PutObjectOutput{ETag: testETag}
	assert.Equal(t, testETag, putOut.ETag)

	now := time.Now().UTC()
	getOut := storage.GetObjectOutput{ContentLength: testLength, LastModified: now}
	assert.Equal(t, int64(testLength), getOut.ContentLength)
	assert.Equal(t, now, getOut.LastModified)
}

func TestPublicModels_Manage(t *testing.T) {
	t.Parallel()

	delIn := storage.DeleteObjectInput{Bucket: testBucket}
	assert.Equal(t, testBucket, delIn.Bucket)

	delOut := storage.DeleteObjectOutput{VersionID: "v1"}
	assert.Equal(t, "v1", delOut.VersionID)

	headOut := storage.HeadObjectOutput{ETag: testETag}
	assert.Equal(t, testETag, headOut.ETag)

	summary := storage.ObjectSummary{Size: testSize}
	assert.Equal(t, int64(testSize), summary.Size)
}

func TestPublicModels_Presign(t *testing.T) {
	t.Parallel()

	presignGet := storage.PresignGetObjectInput{ExpiresIn: time.Hour}
	assert.Equal(t, time.Hour, presignGet.ExpiresIn)

	presignPut := storage.PresignPutObjectInput{ExpiresIn: testExpiry}
	assert.Equal(t, testExpiry, presignPut.ExpiresIn)

	presigned := storage.PresignedURL{URL: testURLString}
	assert.Equal(t, testURLString, presigned.URL)
}

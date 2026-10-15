package contracts_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	testContractImgSize     = 100
	testContractMaxKeys     = 50
	testContractSummarySize = 10
	testContractDuration    = time.Minute
)

func TestModels_Delete(t *testing.T) {
	t.Parallel()

	delIn := contracts.DeleteObjectInput{Bucket: testContractBucket, Key: testContractKey}
	assert.Equal(t, testContractBucket, delIn.Bucket)

	delOut := contracts.DeleteObjectOutput{VersionID: testContractVersionID}
	assert.Equal(t, testContractVersionID, delOut.VersionID)
}

func TestModels_Head(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	headIn := contracts.HeadObjectInput{Bucket: testContractBucket, Key: testContractKey}
	assert.Equal(t, testContractKey, headIn.Key)

	headOut := contracts.HeadObjectOutput{
		LastModified:  now,
		ContentLength: testContractImgSize,
		ContentType:   testContractContentType,
		ETag:          testContractETag,
	}
	assert.Equal(t, int64(testContractImgSize), headOut.ContentLength)
	assert.Equal(t, testContractContentType, headOut.ContentType)
}

func TestModels_List(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	listIn := contracts.ListObjectsInput{
		Bucket:    testContractBucket,
		Prefix:    "p/",
		Delimiter: "/",
		MaxKeys:   testContractMaxKeys,
	}
	assert.Equal(t, int32(testContractMaxKeys), listIn.MaxKeys)

	summary := contracts.ObjectSummary{
		LastModified: now,
		Key:          "p/1",
		ETag:         testContractETag,
		Size:         testContractSummarySize,
	}
	listOut := contracts.ListObjectsOutput{Objects: []contracts.ObjectSummary{summary}, Prefixes: []string{"p/2/"}}
	assert.Len(t, listOut.Objects, 1)
}

func TestModels_Presign(t *testing.T) {
	t.Parallel()

	presignGet := contracts.PresignGetObjectInput{
		Bucket:    testContractBucket,
		Key:       testContractKey,
		ExpiresIn: testContractDuration,
	}
	assert.Equal(t, testContractDuration, presignGet.ExpiresIn)

	presignPut := contracts.PresignPutObjectInput{
		Bucket:      testContractBucket,
		Key:         testContractKey,
		ContentType: testContractContentType,
		ExpiresIn:   testContractDuration,
	}
	assert.Equal(t, testContractContentType, presignPut.ContentType)

	presigned := contracts.PresignedURL{
		URL:       "https://example.com",
		Method:    "PUT",
		ExpiresIn: testContractDuration,
	}
	assert.Equal(t, "https://example.com", presigned.URL)
}

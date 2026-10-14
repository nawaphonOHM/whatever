package contracts_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const (
	testContractBucket      = "my-bucket"
	testContractKey         = "path/to/key.txt"
	testContractContentType = "text/plain"
	testContractETag        = "\"123456\""
	testContractVersionID   = "v1"
	testContractPayloadSize = 9
)

func TestModels_PutInput(t *testing.T) {
	t.Parallel()

	body := bytes.NewReader([]byte("payload"))
	meta := map[string]string{"env": "test"}

	input := contracts.PutObjectInput{
		Bucket:      testContractBucket,
		Key:         testContractKey,
		Body:        body,
		ContentType: testContractContentType,
		Metadata:    meta,
	}
	assert.Equal(t, testContractBucket, input.Bucket)
	assert.Equal(t, testContractKey, input.Key)
	assert.Equal(t, testContractContentType, input.ContentType)
	assert.Equal(t, meta, input.Metadata)
	assert.NotNil(t, input.Body)
}

func TestModels_PutOutput(t *testing.T) {
	t.Parallel()

	output := contracts.PutObjectOutput{
		ETag:      testContractETag,
		VersionID: testContractVersionID,
	}
	assert.Equal(t, testContractETag, output.ETag)
	assert.Equal(t, testContractVersionID, output.VersionID)
}

func TestModels_Get(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	meta := map[string]string{"author": "ohm"}
	rc := io.NopCloser(bytes.NewReader([]byte("test data")))

	getIn := contracts.GetObjectInput{Bucket: testContractBucket, Key: testContractKey}
	assert.Equal(t, testContractBucket, getIn.Bucket)

	getOut := contracts.GetObjectOutput{
		LastModified:  now,
		Body:          rc,
		Metadata:      meta,
		ContentType:   testContractContentType,
		ETag:          testContractETag,
		ContentLength: testContractPayloadSize,
	}
	assert.NotNil(t, getOut.Body)
	assert.Equal(t, int64(testContractPayloadSize), getOut.ContentLength)
	require.NoError(t, getOut.Body.Close())
}

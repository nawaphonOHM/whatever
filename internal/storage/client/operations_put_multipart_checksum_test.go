package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func nonAWSMultipartResponder() func(recordedRequest) *http.Response {
	return func(req recordedRequest) *http.Response {
		if req.Method == http.MethodPut {
			return newHTTPResponse(http.StatusOK, "", map[string]string{headerETag: `"etag"`})
		}
		if req.Query.Has("uploads") {
			return newHTTPResponse(http.StatusOK, initiateXML, nil)
		}
		return newHTTPResponse(http.StatusOK, completeXMLBody, nil)
	}
}

func assertNonAWSMultipartOmissions(t *testing.T, reqs []recordedRequest) {
	t.Helper()
	require.Len(t, reqs, multipartRequestCount)
	assert.Empty(t, reqs[0].Header.Get(headerChecksumAlgo))
	assert.Empty(t, reqs[1].Header.Get(headerChecksumCRC32))
	assert.NotContains(t, string(reqs[completeRequestIndex].Body), "ChecksumCRC32")
}

func runPutObjectMultipartOmissionTest(t *testing.T, endpoint string) {
	t.Helper()
	rec := &recorder{respond: nonAWSMultipartResponder()}
	cfg := testConfigWithKeys()
	cfg.Endpoint = endpoint
	c := connectRecorded(t, cfg, rec)
	out, err := c.PutObject(t.Context(), contracts.PutObjectInput{
		Key: testObjectKey, Body: &zeroReader{left: partSize + 1},
	})
	require.NoError(t, err)
	assert.Equal(t, `"final-etag"`, out.ETag)
	assertNonAWSMultipartOmissions(t, rec.all())
}

func TestPutObject_S3CompatibleOmitsMultipartChecksums(t *testing.T) {
	t.Parallel()
	runPutObjectMultipartOmissionTest(t, "https://minio.internal.net")
}

func TestPutObject_GCSOmitsMultipartChecksums(t *testing.T) {
	t.Parallel()
	runPutObjectMultipartOmissionTest(t, "https://storage.googleapis.com")
}

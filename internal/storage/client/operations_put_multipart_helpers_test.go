package client_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// zeroReader is a non-seekable stream of n zero bytes.
type zeroReader struct{ left int64 }

func (z *zeroReader) Read(p []byte) (int, error) {
	if z.left <= 0 {
		return 0, io.EOF
	}
	n := int(min(int64(len(p)), z.left))
	clear(p[:n])
	z.left -= int64(n)
	return n, nil
}

func multipartResponder(failPart string) func(recordedRequest) *http.Response {
	return func(rec recordedRequest) *http.Response {
		if rec.Method == http.MethodPut {
			return partReply(rec, failPart)
		}
		return controlReply(rec)
	}
}

func partReply(rec recordedRequest, failPart string) *http.Response {
	if rec.Query.Get(partNumberParam) == failPart {
		return newHTTPResponse(http.StatusInternalServerError, "", nil)
	}
	return newHTTPResponse(http.StatusOK, "", map[string]string{"ETag": `"part-etag"`})
}

func controlReply(rec recordedRequest) *http.Response {
	if rec.Method != http.MethodPost {
		return newHTTPResponse(http.StatusNoContent, "", nil)
	}
	if rec.Query.Has("uploads") {
		return newHTTPResponse(http.StatusOK, initiateXML, nil)
	}
	return newHTTPResponse(http.StatusOK, completeXMLBody, nil)
}

type errReader struct{ err error }

func (e *errReader) Read([]byte) (int, error) { return 0, e.err }

func isInitiate(rec recordedRequest) bool {
	return rec.Method == http.MethodPost && rec.Query.Has("uploads")
}

func isComplete(rec recordedRequest) bool {
	return rec.Method == http.MethodPost && !rec.Query.Has("uploads")
}

func failingControl(shouldFail func(recordedRequest) bool) func(recordedRequest) *http.Response {
	ok := multipartResponder("0")
	return func(rec recordedRequest) *http.Response {
		if shouldFail(rec) {
			return xmlError(http.StatusForbidden, "AccessDenied")(rec)
		}
		return ok(rec)
	}
}

func TestPutObject_MultipartInitiateFailureNeedsNoAbort(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: failingControl(isInitiate)}
	_, err := putZeros(t, rec, partSize+1)
	require.ErrorIs(t, err, contracts.ErrAccessDenied)
	assert.Len(t, rec.all(), 1)
}

func TestPutObject_MultipartCompleteFailureAborts(t *testing.T) {
	t.Parallel()

	rec := &recorder{respond: failingControl(isComplete)}
	_, err := putZeros(t, rec, partSize+1)
	require.ErrorIs(t, err, contracts.ErrAccessDenied)
	reqs := rec.all()
	assert.Equal(t, http.MethodDelete, reqs[len(reqs)-1].Method)
}

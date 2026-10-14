package client_test

import (
	"io"
	"net/http"
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

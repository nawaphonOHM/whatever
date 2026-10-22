package rest

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testStreamBytes = 5 * 1024 * 1024

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'A'
	}
	return len(p), nil
}

type countingWriter struct {
	header           http.Header
	committedHeaders http.Header
	code             int
	count            int64
	wroteHeader      bool
}

func (cw *countingWriter) Header() http.Header {
	if cw.header == nil {
		cw.header = make(http.Header)
	}
	return cw.header
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	if !cw.wroteHeader {
		cw.WriteHeader(http.StatusOK)
	}
	cw.count += int64(len(p))
	return len(p), nil
}

func (cw *countingWriter) WriteHeader(code int) {
	if cw.wroteHeader {
		return
	}
	cw.code = code
	cw.wroteHeader = true
	cw.committedHeaders = cw.Header().Clone()
}

func TestStream_StatusCodeDefault(t *testing.T) {
	respZero := Stream(0, nil)
	assert.Equal(t, http.StatusOK, respZero.StatusCode())
	var nilResp *streamResponse
	assert.Equal(t, http.StatusOK, nilResp.StatusCode())
}

func TestStream_LargePayload(t *testing.T) {
	cw := &countingWriter{}
	c, _ := gin.CreateTestContext(cw)
	body := io.NopCloser(io.LimitReader(zeroReader{}, testStreamBytes))
	Stream(http.StatusOK, body).WithHeader("X-Large-Payload", "5MB").Write(c)
	assert.Equal(t, http.StatusOK, cw.code)
	assert.Equal(t, int64(testStreamBytes), cw.count)
	require.NotNil(t, cw.committedHeaders)
	assert.Equal(t, "5MB", cw.committedHeaders.Get("X-Large-Payload"))
}

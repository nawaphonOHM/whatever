package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testStreamDataStr = "data"
	testStreamHdrKey  = "X-Stream-Trace"
	testStreamHdrVal  = "trace-stream-1"
	testStreamChunk   = 4096
)

func TestStream_HeaderChaining(t *testing.T) {
	r := setupTestRouter()
	r.GET("/stream-headers", func(c *gin.Context) {
		body := io.NopCloser(strings.NewReader(testStreamDataStr))
		Stream(http.StatusOK, body).
			WithHeader("x-key", "orig").
			WithHeader("X-Key", "overwritten").
			WithHeaderInt("X-Chunk-Size", testStreamChunk).
			WithHeaderStringer(testStreamHdrKey, customTraceID{id: testStreamHdrVal}).
			Write(c)
	})
	assertChainedHeaders(t, executeGet(t, r, "/stream-headers"))
}

func assertChainedHeaders(t *testing.T, w *httptest.ResponseRecorder) {
	res := w.Result()
	assert.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "overwritten", res.Header.Get("X-Key"))
	assert.Equal(t, "4096", res.Header.Get("X-Chunk-Size"))
	assert.Equal(t, testStreamHdrVal, res.Header.Get(testStreamHdrKey))
	assert.Equal(t, testStreamDataStr, w.Body.String())
}

func TestStream_NilReceiverChaining(t *testing.T) {
	var nilStream *streamResponse
	chained := nilStream.
		WithHeader("X-Test", "val").
		WithHeaderInt("X-Num", testStreamChunk).
		WithHeaderStringer("X-Str", customTraceID{id: "id"})
	assert.Nil(t, chained)
	assert.Equal(t, http.StatusOK, chained.StatusCode())
	c, _ := newTestContext()
	assert.NotPanics(t, func() { chained.Write(c) })
	assert.NotPanics(t, func() { chained.Write(nil) })
}

func TestStream_StringerEdgeCases(t *testing.T) {
	t.Run("nil and typed nil stringer", func(t *testing.T) {
		r := setupTestRouter()
		var typedNil *dummyPointerStringer
		r.GET("/stringer-edge", func(c *gin.Context) {
			Stream(http.StatusOK, io.NopCloser(strings.NewReader(testStreamDataStr))).
				WithHeaderStringer("X-Nil-Interface", nil).
				WithHeaderStringer("X-Typed-Nil", typedNil).
				Write(c)
		})
		res := executeGet(t, r, "/stringer-edge").Result()
		assert.NoError(t, res.Body.Close())
		assert.Empty(t, res.Header.Get("X-Nil-Interface"))
		assert.Empty(t, res.Header.Get("X-Typed-Nil"))
	})
	t.Run("panicking stringer propagates panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Stream(http.StatusOK, nil).WithHeaderStringer("X-Panic", dummyPanicStringer{})
		})
	})
}

func TestStream_EmptyHeaderKey(t *testing.T) {
	r := setupTestRouter()
	r.GET("/empty-header", func(c *gin.Context) {
		Stream(http.StatusOK, io.NopCloser(strings.NewReader(testStreamDataStr))).
			WithHeader("", "ignored").
			WithHeaderInt("", testStreamChunk).
			WithHeaderStringer("", customTraceID{id: "ignored"}).
			WithHeader("X-Empty-Val", "").
			WithHeaderInt("X-Zero", 0).
			Write(c)
	})
	res := executeGet(t, r, "/empty-header").Result()
	assert.NoError(t, res.Body.Close())
	assert.Empty(t, res.Header.Get(""))
	assert.Equal(t, "", res.Header.Get("X-Empty-Val"))
	assert.Equal(t, "0", res.Header.Get("X-Zero"))
}

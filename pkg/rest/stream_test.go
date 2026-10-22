package rest

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testStreamData = "stream payload content"
	hdrStreamType  = "Content-Type"
	hdrStreamText  = "text/plain"
	hdrStreamOctet = "application/octet-stream"
)

func TestStream_OK(t *testing.T) {
	r := setupTestRouter()
	r.GET("/stream-ok", func(c *gin.Context) {
		body := io.NopCloser(strings.NewReader(testStreamData))
		Stream(http.StatusOK, body).WithHeader(hdrStreamType, hdrStreamText).Write(c)
	})
	w := executeGet(t, r, "/stream-ok")
	res := w.Result()
	assert.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, hdrStreamText, res.Header.Get(hdrStreamType))
	assert.Equal(t, testStreamData, w.Body.String())
}

func TestStream_Created_ByteBuffer(t *testing.T) {
	r := setupTestRouter()
	r.POST("/stream-created", func(c *gin.Context) {
		buf := bytes.NewBufferString("created stream")
		Stream(http.StatusCreated, io.NopCloser(buf)).
			WithHeader(hdrStreamType, hdrStreamOctet).
			Write(c)
	})
	w := executePost(t, r, "/stream-created")
	res := w.Result()
	assert.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, hdrStreamOctet, res.Header.Get(hdrStreamType))
	assert.Equal(t, "created stream", w.Body.String())
}

func TestStream_CustomStatus(t *testing.T) {
	r := setupTestRouter()
	r.GET("/stream-partial", func(c *gin.Context) {
		body := io.NopCloser(strings.NewReader(testStreamData))
		Stream(http.StatusPartialContent, body).Write(c)
	})
	w := executeGet(t, r, "/stream-partial")
	res := w.Result()
	assert.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusPartialContent, res.StatusCode)
	assert.Equal(t, testStreamData, w.Body.String())
}

func TestStream_NilBody(t *testing.T) {
	r := setupTestRouter()
	r.GET("/stream-nil", func(c *gin.Context) {
		Stream(http.StatusNoContent, nil).
			WithHeader("X-Nil-Stream", "active").
			Write(c)
	})
	w := executeGet(t, r, "/stream-nil")
	res := w.Result()
	assert.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, "active", res.Header.Get("X-Nil-Stream"))
	assert.Empty(t, w.Body.String())
}

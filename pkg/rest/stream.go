package rest

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// streamResponse represents an HTTP response that streams content from an io.ReadCloser.
type streamResponse struct {
	headers    map[string]string
	body       io.ReadCloser
	statusCode int
}

// Stream creates a streaming Response with the given status code and body reader.
// When Write is invoked, headers and status code are written, the body is streamed
// directly to the client, and body.Close() is guaranteed to execute.
//
// Ownership Note: If Stream is called but Write is never executed (e.g., if the
// handler aborts early and returns a different response), the caller retains
// responsibility for closing the body to prevent resource leaks.
func Stream(statusCode int, body io.ReadCloser) Response {
	return &streamResponse{
		statusCode: statusCode,
		body:       body,
	}
}

// StatusCode returns the configured HTTP status code, defaulting to 200 OK
// if unset or if the receiver is nil.
func (r *streamResponse) StatusCode() int {
	if r == nil || r.statusCode == 0 {
		return http.StatusOK
	}
	return r.statusCode
}

// Write renders the streaming response to the given Gin context.
// It guarantees that body.Close() is called whenever body is non-nil,
// even if the context is nil or an error occurs during transfer.
func (r *streamResponse) Write(c *gin.Context) {
	if r == nil {
		return
	}
	defer closeStream(r.body, c)
	r.writeTo(c)
}

func (r *streamResponse) writeTo(c *gin.Context) {
	if c == nil || c.Writer == nil {
		return
	}
	contracts.WriteHeaders(c, r.headers)
	c.Status(r.StatusCode())
	copyStreamBody(c, r.body)
}

func copyStreamBody(c *gin.Context, body io.Reader) {
	if body == nil {
		return
	}
	if _, err := io.Copy(c.Writer, body); err != nil {
		recordError(c, err)
	}
}

func closeStream(body io.Closer, c *gin.Context) {
	if body == nil {
		return
	}
	if err := body.Close(); err != nil {
		recordError(c, err)
	}
}

func recordError(c *gin.Context, err error) {
	if c == nil {
		return
	}
	if ginErr := c.Error(err); ginErr != nil {
		return
	}
}

// WithHeader attaches a custom key-value header to the streaming response.
func (r *streamResponse) WithHeader(key, value string) Response {
	if r == nil {
		return r
	}
	contracts.SetHeader(&r.headers, key, value)
	return r
}

// WithHeaderInt attaches a custom header whose integer value is converted to a string.
func (r *streamResponse) WithHeaderInt(key string, value int) Response {
	if r == nil {
		return r
	}
	return r.WithHeader(key, strconv.Itoa(value))
}

// WithHeaderStringer attaches a custom header whose value is obtained from fmt.Stringer.
func (r *streamResponse) WithHeaderStringer(key string, value fmt.Stringer) Response {
	if r == nil {
		return r
	}
	return applyStreamHeaderStringer(r, key, value)
}

func applyStreamHeaderStringer(r *streamResponse, key string, value fmt.Stringer) Response {
	str, ok := contracts.SafeStringer(value)
	if !ok {
		return r
	}
	return r.WithHeader(key, str)
}

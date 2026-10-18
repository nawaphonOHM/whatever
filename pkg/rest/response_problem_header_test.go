package rest

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testProbResourceNum  = 404
	testProbRetryAfter   = 120
	testHeaderRetryAfter = "X-Retry-After"
)

func TestResponseHeaders_ProblemResponses(t *testing.T) {
	t.Run("NotFound with custom headers and stringer", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/not-found", func(c *gin.Context) {
			NotFound("NOT_FOUND", "resource missing").
				WithHeader("X-Resource-ID", "item-404").
				WithHeaderInt("X-Resource-Num", testProbResourceNum).
				WithHeaderStringer(testTraceHeader, customTraceID{id: "trace-404"}).
				WithHeader("content-type", "text/plain").
				Write(c)
		})

		w := executeGet(t, r, "/not-found")
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
		assert.Equal(t, "item-404", w.Header().Get("X-Resource-ID"))
		assert.Equal(t, "404", w.Header().Get("X-Resource-Num"))
		assert.Equal(t, "trace-404", w.Header().Get(testTraceHeader))
	})

	t.Run("BadRequest with custom headers", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/bad-request", func(c *gin.Context) {
			BadRequest("BAD_PARAM", "invalid param").
				WithHeader("X-Error-Reason", "invalid_param").
				WithHeaderInt(testHeaderRetryAfter, testProbRetryAfter).
				Write(c)
		})

		w := executeGet(t, r, "/bad-request")
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
		assert.Equal(t, "invalid_param", w.Header().Get("X-Error-Reason"))
		assert.Equal(t, "120", w.Header().Get(testHeaderRetryAfter))
	})

	testOtherProblemResponsesWithHeaders(t)
}

func testOtherProblemResponsesWithHeaders(t *testing.T) {
	t.Run("Auth and Server Error problem helpers with headers", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/unauth", func(c *gin.Context) {
			Unauthorized("UNAUTHORIZED", "unauthorized").
				WithHeader("WWW-Authenticate", "Bearer").
				Write(c)
		})
		r.GET("/server-err", func(c *gin.Context) {
			InternalServerError("INTERNAL", "internal").
				WithHeader(testHeaderRetryAfter, "120").
				Write(c)
		})

		wUnauth := executeGet(t, r, "/unauth")
		assert.Equal(t, http.StatusUnauthorized, wUnauth.Code)
		assert.Equal(t, "Bearer", wUnauth.Header().Get("WWW-Authenticate"))

		wErr := executeGet(t, r, "/server-err")
		assert.Equal(t, http.StatusInternalServerError, wErr.Code)
		assert.Equal(t, "120", wErr.Header().Get(testHeaderRetryAfter))
	})

	t.Run("Forbidden and custom error helpers with headers", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/forbidden", func(c *gin.Context) {
			Forbidden("FORBIDDEN", "forbidden").
				WithHeader("X-Scope-Required", "admin").
				Write(c)
		})
		r.GET("/custom-err", func(c *gin.Context) {
			Error(http.StatusTeapot, "TEAPOT", "teapot").
				WithHeader("X-Teapot", "true").
				Write(c)
		})

		wForb := executeGet(t, r, "/forbidden")
		assert.Equal(t, http.StatusForbidden, wForb.Code)
		assert.Equal(t, "admin", wForb.Header().Get("X-Scope-Required"))

		wTeapot := executeGet(t, r, "/custom-err")
		assert.Equal(t, http.StatusTeapot, wTeapot.Code)
		assert.Equal(t, "true", wTeapot.Header().Get("X-Teapot"))
	})
}

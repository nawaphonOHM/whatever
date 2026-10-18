package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTraceHeader = "X-Trace-ID"
	testCreatedNum  = 42
)

type customTraceID struct{ id string }

func (t customTraceID) String() string { return t.id }

func TestResponseHeaders_SuccessResponses(t *testing.T) {
	t.Run("OK with headers and stringer", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/ok", func(c *gin.Context) {
			OK("data-val").
				WithHeader("X-String-Header", "hello").
				WithHeaderStringer(testTraceHeader, customTraceID{id: "trace-123"}).
				Write(c)
		})

		w := executeGet(t, r, "/ok")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "hello", w.Header().Get("X-String-Header"))
		assert.Equal(t, "trace-123", w.Header().Get(testTraceHeader))

		var body SuccessResponse[string]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.Equal(t, "data-val", body.Data)
		assert.True(t, body.Success)
	})

	t.Run("Created with headers", func(t *testing.T) {
		r := setupTestRouter()
		r.POST("/created", func(c *gin.Context) {
			Created(map[string]int{"id": testCreatedNum}).
				WithHeader("Location", "/items/42").
				Write(c)
		})

		w := executePost(t, r, "/created")
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "/items/42", w.Header().Get("Location"))
	})

	t.Run("JSON custom status with headers", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/accepted", func(c *gin.Context) {
			JSON(http.StatusAccepted, "accepted").
				WithHeader("X-Task-ID", "task-99").
				Write(c)
		})

		w := executeGet(t, r, "/accepted")
		assert.Equal(t, http.StatusAccepted, w.Code)
		assert.Equal(t, "task-99", w.Header().Get("X-Task-ID"))
	})
}

func TestResponseHeaders_NoContent(t *testing.T) {
	t.Run("NoContent with header and stringer", func(t *testing.T) {
		r := setupTestRouter()
		r.DELETE("/no-content", func(c *gin.Context) {
			NoContent().
				WithHeader("X-Action", "deleted").
				WithHeaderStringer(testTraceHeader, customTraceID{id: "del-77"}).
				Write(c)
		})

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/no-content", nil)
		require.NoError(t, err)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "deleted", w.Header().Get("X-Action"))
		assert.Equal(t, "del-77", w.Header().Get(testTraceHeader))
		assert.Empty(t, w.Body.String())
	})
}

func TestResponseHeaders_ChainingAndOverwrite(t *testing.T) {
	t.Run("overwrites duplicate header keys with mixed casing", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/overwrite", func(c *gin.Context) {
			OK("test").
				WithHeader("x-key", "initial").
				WithHeader("X-Other", "other").
				WithHeader("X-Key", "overwritten-1").
				WithHeader("X-KEY", "overwritten-2").
				Write(c)
		})

		w := executeGet(t, r, "/overwrite")
		assert.Equal(t, "overwritten-2", w.Header().Get("X-Key"))
		assert.Equal(t, "other", w.Header().Get("X-Other"))
	})
}

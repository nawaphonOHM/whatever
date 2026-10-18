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
	testIntHeaderCount = "X-Total-Count"
	testIntHeaderItem  = "X-Item-ID"
	testIntHeaderTask  = "X-Task-Count"
	testIntValCount    = 42
	testIntValItem     = 101
	testIntValDeleted  = 5
	testIntValNeg      = -50
	testIntValOver1    = 20
	testIntValOver2    = 100
)

func TestResponseHeaders_WithHeaderInt_Success(t *testing.T) {
	t.Run("OK with WithHeaderInt", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/ok-int", func(c *gin.Context) {
			OK("data-val").WithHeaderInt(testIntHeaderCount, testIntValCount).Write(c)
		})

		w := executeGet(t, r, "/ok-int")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "42", w.Header().Get(testIntHeaderCount))

		var body SuccessResponse[string]
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		assert.True(t, body.Success)
	})

	t.Run("Created and JSON with WithHeaderInt", func(t *testing.T) {
		r := setupTestRouter()
		r.POST("/created-int", func(c *gin.Context) {
			Created("item").WithHeaderInt(testIntHeaderItem, testIntValItem).Write(c)
		})
		r.GET("/json-int", func(c *gin.Context) {
			JSON(http.StatusAccepted, "task").WithHeaderInt(testIntHeaderTask, 0).Write(c)
		})

		wCreated := executePost(t, r, "/created-int")
		assert.Equal(t, http.StatusCreated, wCreated.Code)
		assert.Equal(t, "101", wCreated.Header().Get(testIntHeaderItem))

		wJSON := executeGet(t, r, "/json-int")
		assert.Equal(t, http.StatusAccepted, wJSON.Code)
		assert.Equal(t, "0", wJSON.Header().Get(testIntHeaderTask))
	})
}

func TestResponseHeaders_WithHeaderInt_NoContentAndOverwrite(t *testing.T) {
	t.Run("NoContent with WithHeaderInt", func(t *testing.T) {
		r := setupTestRouter()
		r.DELETE("/no-content-int", func(c *gin.Context) {
			NoContent().WithHeaderInt("X-Deleted-Count", testIntValDeleted).Write(c)
		})

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/no-content-int", nil)
		require.NoError(t, err)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "5", w.Header().Get("X-Deleted-Count"))
	})

	t.Run("overwrites with mixed types and case", func(t *testing.T) {
		r := setupTestRouter()
		r.GET("/overwrite-int", func(c *gin.Context) {
			OK("test").
				WithHeader("x-key", "initial").
				WithHeaderInt("X-Key", testIntValOver1).
				WithHeader("X-KEY", "overwritten-str").
				WithHeaderInt("X-A", testIntValNeg).
				WithHeaderInt("x-a", testIntValOver2).
				Write(c)
		})

		w := executeGet(t, r, "/overwrite-int")
		assert.Equal(t, "overwritten-str", w.Header().Get("X-Key"))
		assert.Equal(t, "100", w.Header().Get("X-A"))
	})
}

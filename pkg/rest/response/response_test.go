package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Constants for test fixture values.
const (
	testItemKey   = "item"
	testItemVal   = "widget"
	testMsgOK     = "Fetched successfully"
	testCreatedID = 42
	testMsgCreate = "Resource created"
)

// init initializes the gin engine in test mode.
func init() {
	gin.SetMode(gin.TestMode)
}

// setupTestRouter creates a new gin engine for test cases.
func setupTestRouter() *gin.Engine {
	return gin.New()
}

// executeGet executes a test GET request against the router.
// Returns the response recorder for subsequent assertions.
func executeGet(
	t *testing.T,
	r *gin.Engine,
	path string,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)
	return w
}

// executePost executes a test POST request against the router.
// Returns the response recorder for subsequent assertions.
func executePost(
	t *testing.T,
	r *gin.Engine,
	path string,
) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, path, nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)
	return w
}

// TestOK tests standard 200 OK responses with generic payload.
func TestOK(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test-ok", func(c *gin.Context) {
		OK(map[string]string{testItemKey: testItemVal}).Write(c)
	})

	w := executeGet(t, r, "/test-ok")
	assert.Equal(t, http.StatusOK, w.Code)

	var resp SuccessResponse[map[string]string]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp.Success)
	assert.Equal(t, testItemVal, resp.Data[testItemKey])
}

// TestOK_WithMessage tests 200 OK with a custom message string.
func TestOK_WithMessage(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test-ok-message", func(c *gin.Context) {
		OK([]string{"item1", "item2"}, testMsgOK).Write(c)
	})

	w := executeGet(t, r, "/test-ok-message")
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	val, ok := resp["success"].(bool)
	require.True(t, ok)
	assert.True(t, val)
	assert.Equal(t, testMsgOK, resp["message"])
}

// TestCreated tests 201 Created responses with payload and message.
func TestCreated(t *testing.T) {
	r := setupTestRouter()
	r.POST("/test-created", func(c *gin.Context) {
		Created(map[string]int{"id": testCreatedID}, testMsgCreate).Write(c)
	})

	w := executePost(t, r, "/test-created")
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp SuccessResponse[map[string]int]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp.Success)
	assert.Equal(t, testCreatedID, resp.Data["id"])
}

// TestNoContent tests 204 No Content responses with empty body.
func TestNoContent(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/test-no-content", func(c *gin.Context) {
		NoContent().Write(c)
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/test-no-content", nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

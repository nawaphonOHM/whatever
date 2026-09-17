package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	r := gin.New()
	return r
}

func TestOK(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test-ok", func(c *gin.Context) {
		OK(map[string]string{"item": "widget"}).Write(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-ok", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SuccessResponse[map[string]string]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Equal(t, "widget", resp.Data["item"])
	assert.Empty(t, resp.Message)
	assert.False(t, resp.Timestamp.IsZero())
}

func TestOK_WithMessage(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test-ok-message", func(c *gin.Context) {
		OK([]string{"item1", "item2"}, "Fetched successfully").Write(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test-ok-message", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "Fetched successfully", resp["message"])
	assert.NotNil(t, resp["data"])
	assert.NotEmpty(t, resp["timestamp"])
}

func TestCreated(t *testing.T) {
	r := setupTestRouter()
	r.POST("/test-created", func(c *gin.Context) {
		Created(map[string]int{"id": 42}, "Resource created").Write(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test-created", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp SuccessResponse[map[string]int]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Equal(t, 42, resp.Data["id"])
	assert.Equal(t, "Resource created", resp.Message)
	assert.False(t, resp.Timestamp.IsZero())
}

func TestNoContent(t *testing.T) {
	r := setupTestRouter()
	r.DELETE("/test-no-content", func(c *gin.Context) {
		NoContent().Write(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/test-no-content", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		handler        gin.HandlerFunc
		expectedStatus int
		expectedTitle  string
		expectedCode   string
		expectedDetail string
	}{
		{
			name:     "BadRequest",
			endpoint: "/bad-request",
			handler: func(c *gin.Context) {
				BadRequest("INVALID_INPUT", "Input validation failed", map[string]string{"field": "email"}).Write(c)
			},
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "INVALID_INPUT",
			expectedDetail: "Input validation failed",
		},
		{
			name:     "Unauthorized",
			endpoint: "/unauthorized",
			handler: func(c *gin.Context) {
				Unauthorized("UNAUTHORIZED", "Missing token").Write(c)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedTitle:  "Unauthorized",
			expectedCode:   "UNAUTHORIZED",
			expectedDetail: "Missing token",
		},
		{
			name:     "Forbidden",
			endpoint: "/forbidden",
			handler: func(c *gin.Context) {
				Forbidden("FORBIDDEN", "Access denied").Write(c)
			},
			expectedStatus: http.StatusForbidden,
			expectedTitle:  "Forbidden",
			expectedCode:   "FORBIDDEN",
			expectedDetail: "Access denied",
		},
		{
			name:     "NotFound",
			endpoint: "/not-found",
			handler: func(c *gin.Context) {
				NotFound("RESOURCE_NOT_FOUND", "User not found").Write(c)
			},
			expectedStatus: http.StatusNotFound,
			expectedTitle:  "Not Found",
			expectedCode:   "RESOURCE_NOT_FOUND",
			expectedDetail: "User not found",
		},
		{
			name:     "InternalServerError",
			endpoint: "/internal-error",
			handler: func(c *gin.Context) {
				InternalServerError("INTERNAL_ERROR", "Unexpected server error").Write(c)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedTitle:  "Internal Server Error",
			expectedCode:   "INTERNAL_ERROR",
			expectedDetail: "Unexpected server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.GET(tt.endpoint, tt.handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.endpoint, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

			var prob ProblemDetails
			err := json.Unmarshal(w.Body.Bytes(), &prob)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, prob.Status)
			assert.Equal(t, "about:blank", prob.Type)
			assert.Equal(t, tt.expectedTitle, prob.Title)
			assert.Equal(t, tt.expectedCode, prob.Code)
			assert.Equal(t, tt.expectedDetail, prob.Detail)
			assert.Equal(t, tt.endpoint, prob.Instance)
		})
	}
}

func TestNewProblemDetails_MultipleDetails(t *testing.T) {
	prob := NewProblemDetails(http.StatusBadRequest, "VALIDATION_ERROR", "Multiple issues found", "issue 1", "issue 2")
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, "about:blank", prob.Type)
	assert.Equal(t, "Bad Request", prob.Title)
	assert.Equal(t, "VALIDATION_ERROR", prob.Code)
	assert.Equal(t, "Multiple issues found", prob.Detail)
	assert.Equal(t, []interface{}{"issue 1", "issue 2"}, prob.Details)
}

func TestSuccessResponse_ShapeSerialization(t *testing.T) {
	fixedTime := time.Date(2026, 9, 17, 2, 30, 0, 0, time.UTC)
	resp := NewSuccessResponse("payload", "operation succeeded").WithTimestamp(fixedTime)

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	require.NoError(t, err)

	assert.Equal(t, true, raw["success"])
	assert.Equal(t, "operation succeeded", raw["message"])
	assert.Equal(t, "payload", raw["data"])
	assert.Equal(t, "2026-09-17T02:30:00Z", raw["timestamp"])

	// Test with chained message
	resp2 := NewSuccessResponse(123).WithMessage("custom").WithTimestamp(fixedTime)
	assert.Equal(t, "custom", resp2.Message)
	assert.Equal(t, 123, resp2.Data)
	assert.Equal(t, fixedTime, resp2.Timestamp)
	assert.True(t, resp2.Success)
}

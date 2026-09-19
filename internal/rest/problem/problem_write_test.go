package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func assertDecodedFields(t *testing.T, body []byte) {
	t.Helper()
	var decoded ProblemDetails
	require.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, http.StatusNotFound, decoded.Status)
	assert.Equal(t, DefaultProblemType, decoded.Type)
	assert.Equal(t, "Not Found", decoded.Title)
	assert.Equal(t, "Resource was not located", decoded.Detail)
	assert.Equal(t, "/test/path", decoded.Instance)
}

func TestProblemDetails_WriteDefaults(t *testing.T) {
	gc, w := newTestContext()
	gc.Request = httptest.NewRequest(http.MethodGet, "/test/path", nil)

	prob := &ProblemDetails{
		Status: http.StatusNotFound,
		Detail: "Resource was not located",
	}
	prob.Write(gc)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
	assertDecodedFields(t, w.Body.Bytes())
}

func TestProblemDetails_WriteCustomFields(t *testing.T) {
	gc, w := newTestContext()
	gc.Request = httptest.NewRequest(http.MethodGet, "/orig", nil)

	prob := &ProblemDetails{
		Type:     "https://example.com/err",
		Title:    "Custom Title",
		Status:   http.StatusUnprocessableEntity,
		Detail:   "Custom detail",
		Instance: "/custom/instance",
		Code:     "CUSTOM_CODE",
	}
	prob.Write(gc)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "https://example.com/err", prob.Type)
	assert.Equal(t, "Custom Title", prob.Title)
	assert.Equal(t, "/custom/instance", prob.Instance)
}

func TestProblemDetails_WriteNilChecks(t *testing.T) {
	gc, _ := newTestContext()
	var nilProb *ProblemDetails
	nilProb.Write(gc)

	prob := &ProblemDetails{Status: http.StatusBadRequest}
	prob.Write(nil)

	gcNilReq, _ := newTestContext()
	gcNilReq.Request = nil
	prob.Write(gcNilReq)
	assert.Empty(t, prob.Instance)
}

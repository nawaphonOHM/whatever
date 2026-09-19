package problem

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testUnknownStatus = 999
	testCustomCode    = "CUSTOM_ERR"
	testValidCode     = "BAD_REQ"
	testMultiCode     = "VALIDATION_ERROR"
	testNotFoundCode  = "NOT_FOUND"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNew_SingleDetail(t *testing.T) {
	prob := New(
		http.StatusBadRequest,
		testValidCode,
		"Invalid field",
		"field 'email' is invalid",
	)
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, DefaultProblemType, prob.Type)
	assert.Equal(t, "Bad Request", prob.Title)
	assert.Equal(t, testValidCode, prob.Code)
	assert.Equal(t, "Invalid field", prob.Detail)
	assert.Equal(t, "field 'email' is invalid", prob.Details)
}

func TestNew_MultipleDetails(t *testing.T) {
	prob := New(
		http.StatusBadRequest,
		testMultiCode,
		"Multiple issues",
		"issue 1",
		"issue 2",
	)
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, DefaultProblemType, prob.Type)
	assert.Equal(t, "Bad Request", prob.Title)
	assert.Equal(t, testMultiCode, prob.Code)
	assert.Equal(t, "Multiple issues", prob.Detail)
	assert.Equal(t, []any{"issue 1", "issue 2"}, prob.Details)
}

func TestNew_NoDetails(t *testing.T) {
	prob := New(http.StatusNotFound, testNotFoundCode, "Item not found")
	assert.Equal(t, http.StatusNotFound, prob.Status)
	assert.Equal(t, DefaultProblemType, prob.Type)
	assert.Equal(t, "Not Found", prob.Title)
	assert.Equal(t, testNotFoundCode, prob.Code)
	assert.Equal(t, "Item not found", prob.Detail)
	assert.Nil(t, prob.Details)
}

func TestProblemDetails_StatusCode(t *testing.T) {
	var nilProb *ProblemDetails
	assert.Equal(t, http.StatusInternalServerError, nilProb.StatusCode())

	zeroProb := &ProblemDetails{}
	assert.Equal(t, http.StatusInternalServerError, zeroProb.StatusCode())

	customProb := &ProblemDetails{Status: http.StatusTeapot}
	assert.Equal(t, http.StatusTeapot, customProb.StatusCode())
}

func TestProblemDetails_UnknownStatusTitle(t *testing.T) {
	probCustom := New(testUnknownStatus, testCustomCode, "Unknown error status")
	assert.Equal(t, "Error", probCustom.Title)

	probZero := &ProblemDetails{}
	probZero.ensureTitle()
	assert.Equal(t, "Internal Server Error", probZero.Title)
}

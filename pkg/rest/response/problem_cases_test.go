package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// problemTestCase defines table-driven test cases for RFC 9457 Problem Details.
type problemTestCase struct {
	handler        gin.HandlerFunc
	name           string
	endpoint       string
	expectedTitle  string
	expectedCode   string
	expectedDetail string
	expectedStatus int
}

// buildProblemCases constructs problem test fixture cases.
func buildProblemCases() []problemTestCase {
	return []problemTestCase{
		{
			handler: func(c *gin.Context) {
				BadRequest("INVALID_INPUT", "Input error", "field").Write(c)
			},
			name:           "BadRequest",
			endpoint:       "/bad-request",
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "INVALID_INPUT",
			expectedDetail: "Input error",
		},
		{
			handler: func(c *gin.Context) {
				Unauthorized("UNAUTHORIZED", "Missing token").Write(c)
			},
			name:           "Unauthorized",
			endpoint:       "/unauthorized",
			expectedStatus: http.StatusUnauthorized,
			expectedTitle:  "Unauthorized",
			expectedCode:   "UNAUTHORIZED",
			expectedDetail: "Missing token",
		},
		{
			handler: func(c *gin.Context) {
				Forbidden("FORBIDDEN", "Access denied").Write(c)
			},
			name:           "Forbidden",
			endpoint:       "/forbidden",
			expectedStatus: http.StatusForbidden,
			expectedTitle:  "Forbidden",
			expectedCode:   "FORBIDDEN",
			expectedDetail: "Access denied",
		},
		{
			handler: func(c *gin.Context) {
				NotFound("NOT_FOUND", "User not found").Write(c)
			},
			name:           "NotFound",
			endpoint:       "/not-found",
			expectedStatus: http.StatusNotFound,
			expectedTitle:  "Not Found",
			expectedCode:   "NOT_FOUND",
			expectedDetail: "User not found",
		},
		{
			handler: func(c *gin.Context) {
				InternalServerError("ERR", "Server error").Write(c)
			},
			name:           "InternalServerError",
			endpoint:       "/internal-error",
			expectedStatus: http.StatusInternalServerError,
			expectedTitle:  "Internal Server Error",
			expectedCode:   "ERR",
			expectedDetail: "Server error",
		},
	}
}

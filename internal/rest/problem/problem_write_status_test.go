package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPProblemDetails_WriteExplicitStatuses(t *testing.T) {
	testCases := []struct {
		name          string
		expectedTitle string
		status        int
	}{
		{
			name:          "bad request",
			expectedTitle: "Bad Request",
			status:        http.StatusBadRequest,
		},
		{
			name:          "unauthorized",
			expectedTitle: "Unauthorized",
			status:        http.StatusUnauthorized,
		},
		{
			name:          "forbidden",
			expectedTitle: "Forbidden",
			status:        http.StatusForbidden,
		},
		{
			name:          "not found",
			expectedTitle: "Not Found",
			status:        http.StatusNotFound,
		},
		{
			name:          "unprocessable entity",
			expectedTitle: "Unprocessable Entity",
			status:        http.StatusUnprocessableEntity,
		},
		{
			name:          "service unavailable",
			expectedTitle: "Service Unavailable",
			status:        http.StatusServiceUnavailable,
		},
		{
			name:          "teapot",
			expectedTitle: "I'm a teapot",
			status:        http.StatusTeapot,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gc, w := newTestContext()
			gc.Request = httptest.NewRequest(http.MethodGet, "/explicit-status", nil)

			prob := &PProblemDetails{
				Status: tc.status,
				Detail: "explicit status detail",
			}
			prob.Write(gc)

			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, MediaTypeProblemJSON, w.Header().Get(headerContentType))
			assert.Equal(t, tc.status, prob.Status)

			var decoded PProblemDetails
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &decoded))
			assert.Equal(t, tc.status, decoded.Status)
			assert.Equal(t, tc.expectedTitle, decoded.Title)
			assert.Equal(t, DefaultProblemType, decoded.Type)
			assert.Equal(t, "explicit status detail", decoded.Detail)
			assert.Equal(t, "/explicit-status", decoded.Instance)
		})
	}
}

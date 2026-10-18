package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPProblemDetails_WriteZeroStatus(t *testing.T) {
	t.Run("empty problem struct defaults status to 500", func(t *testing.T) {
		gc, w := newTestContext()
		gc.Request = httptest.NewRequest(http.MethodGet, "/zero-status", nil)

		prob := &PProblemDetails{}
		prob.Write(gc)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get(headerContentType))
		assert.Equal(t, http.StatusInternalServerError, prob.Status)

		var decoded PProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &decoded))
		assert.Equal(t, http.StatusInternalServerError, decoded.Status)
		assert.Equal(t, DefaultProblemType, decoded.Type)
		assert.Equal(t, "Internal Server Error", decoded.Title)
		assert.Equal(t, "/zero-status", decoded.Instance)
	})

	t.Run("problem with detail but zero status defaults status to 500", func(t *testing.T) {
		gc, w := newTestContext()
		gc.Request = httptest.NewRequest(http.MethodPost, "/zero-status-detail", nil)

		prob := &PProblemDetails{
			Detail: "unexpected crash",
		}
		prob.Write(gc)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get(headerContentType))
		assert.Equal(t, http.StatusInternalServerError, prob.Status)

		var decoded PProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &decoded))
		assert.Equal(t, http.StatusInternalServerError, decoded.Status)
		assert.Equal(t, DefaultProblemType, decoded.Type)
		assert.Equal(t, "Internal Server Error", decoded.Title)
		assert.Equal(t, "unexpected crash", decoded.Detail)
		assert.Equal(t, "/zero-status-detail", decoded.Instance)
	})
}

package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProbErrNum = 1001
	testProbNum3   = 3
	testProbNum4   = 4
	testProbNum42  = 42
	testProbNum500 = 500
)

type probTraceID struct{ id string }

func (p probTraceID) String() string { return p.id }

type probPanicStringer struct{}

func (probPanicStringer) String() string { panic("prob panic") }

func TestPProblemDetails_WithHeader(t *testing.T) {
	t.Run("sets custom headers and preserves content type", func(t *testing.T) {
		gc, w := newTestContext()
		gc.Request = httptest.NewRequest(http.MethodGet, "/err-path", nil)

		prob := New(http.StatusBadRequest, "BAD_REQ", "Invalid input").
			WithHeader("X-Error-Code", "1001").
			WithHeaderInt("X-Error-Num", testProbErrNum).
			WithHeaderStringer("X-Trace-ID", probTraceID{id: "trace-999"})

		prob.Write(gc)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
		assert.Equal(t, "1001", w.Header().Get("X-Error-Code"))
		assert.Equal(t, "1001", w.Header().Get("X-Error-Num"))
		assert.Equal(t, "trace-999", w.Header().Get("X-Trace-ID"))

		var decoded PProblemDetails
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &decoded))
		assert.Equal(t, "BAD_REQ", decoded.Code)
		assert.Equal(t, "Invalid input", decoded.Detail)
	})

	t.Run("content type header override protection", func(t *testing.T) {
		gc, w := newTestContext()
		gc.Request = httptest.NewRequest(http.MethodGet, "/err-path", nil)

		prob := New(http.StatusBadRequest, "BAD_REQ", "Invalid input").
			WithHeader("Content-Type", "text/plain").
			WithHeader("content-type", "application/xml")

		prob.Write(gc)

		assert.Equal(t, MediaTypeProblemJSON, w.Header().Get("Content-Type"))
	})

	t.Run("overwrites duplicate headers with mixed case and ints", func(t *testing.T) {
		gc, w := newTestContext()
		prob := New(http.StatusNotFound, "NOT_FOUND", "not found").
			WithHeader("x-a", "1").
			WithHeader("X-A", "2").
			WithHeaderInt("x-A", testProbNum3).
			WithHeaderInt("X-A", testProbNum4)

		prob.Write(gc)
		assert.Equal(t, "4", w.Header().Get("X-A"))
	})
}

func TestPProblemDetails_HeaderEdgeCases(t *testing.T) {
	t.Run("nil receiver chaining", func(t *testing.T) {
		gc, _ := newTestContext()
		var nilProb *PProblemDetails
		chainedProb := nilProb.
			WithHeader("X-K", "V").
			WithHeaderInt("X-Count", testProbNum42).
			WithHeaderStringer("X-K", probTraceID{id: "t"}).
			WithHeader("X-Other", "other")
		assert.Nil(t, chainedProb)
		assert.Equal(t, http.StatusInternalServerError, chainedProb.StatusCode())
		assert.NotPanics(t, func() { chainedProb.Write(gc) })
		assert.NotPanics(t, func() { chainedProb.Write(nil) })
	})

	t.Run("nil stringers and empty keys", func(t *testing.T) {
		gc, w := newTestContext()
		var typedNil *probTraceID
		prob := New(http.StatusInternalServerError, "INTERNAL", "server error").
			WithHeaderStringer("X-Nil", nil).
			WithHeaderStringer("X-TypedNil", typedNil).
			WithHeader("", "ignored-empty-key").
			WithHeaderInt("", testProbNum500).
			WithHeader("X-Empty-Val", "").
			WithHeaderInt("X-Count", 0)

		prob.Write(gc)

		assert.Empty(t, w.Header().Get("X-Nil"))
		assert.Empty(t, w.Header().Get("X-TypedNil"))
		assert.NotContains(t, w.Header(), "")
		assert.Contains(t, w.Header(), "X-Empty-Val")
		assert.Equal(t, "0", w.Header().Get("X-Count"))
	})

	t.Run("panicking stringer propagates panic", func(t *testing.T) {
		assert.Panics(t, func() {
			New(http.StatusInternalServerError, "INTERNAL", "server error").
				WithHeaderStringer("X-Panic", probPanicStringer{})
		})
	})
}

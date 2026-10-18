package rest

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testDefaultKey  = "X-Key"
	testPayloadData = "data"
	testEdgeNum42   = 42
	testEdgeNum100  = 100
	testEdgeNumNeg  = -50
	testEdgeNum10   = 10
)

type dummyPointerStringer struct{ val string }

func (p *dummyPointerStringer) String() string { return p.val }

type dummyPanicStringer struct{}

func (dummyPanicStringer) String() string { panic("stringer panic") }

func TestResponseHeaders_NilReceivers(t *testing.T) {
	c, _ := newTestContext()
	testNilSuccessChaining(t, c)
	testNilNoContentChaining(t, c)
}

func testNilSuccessChaining(t *testing.T, c *gin.Context) {
	var nilSuccess *successResponse[string]
	chained := nilSuccess.
		WithHeader(testDefaultKey, "value").
		WithHeaderInt("X-Count", testEdgeNum42).
		WithHeaderStringer(testDefaultKey, &dummyPointerStringer{val: "v"}).
		WithHeader("X-Other", "other")
	assert.Nil(t, chained)
	assert.Equal(t, http.StatusOK, chained.StatusCode())
	assert.NotPanics(t, func() { chained.Write(c) })
	assert.NotPanics(t, func() { chained.Write(nil) })
}

func testNilNoContentChaining(t *testing.T, c *gin.Context) {
	var nilNoContent *noContentResponse
	chained := nilNoContent.
		WithHeader(testDefaultKey, "value").
		WithHeaderInt("X-Count", testEdgeNum42).
		WithHeaderStringer(testDefaultKey, &dummyPointerStringer{val: "v"}).
		WithHeader("X-Other", "other")
	assert.Nil(t, chained)
	assert.Equal(t, http.StatusNoContent, chained.StatusCode())
	assert.NotPanics(t, func() { chained.Write(c) })
	assert.NotPanics(t, func() { chained.Write(nil) })
}

func TestResponseHeaders_StringerEdgeCases(t *testing.T) {
	t.Run("nil and typed nil stringer", func(t *testing.T) {
		c, w := newTestContext()
		var typedNil *dummyPointerStringer
		resp := OK(testPayloadData).
			WithHeaderStringer("X-Nil-Interface", nil).
			WithHeaderStringer("X-Typed-Nil", typedNil)

		resp.Write(c)

		assert.Empty(t, w.Header().Get("X-Nil-Interface"))
		assert.Empty(t, w.Header().Get("X-Typed-Nil"))
	})

	t.Run("panicking stringer propagates panic", func(t *testing.T) {
		assert.Panics(t, func() {
			OK(testPayloadData).WithHeaderStringer("X-Panicking", dummyPanicStringer{})
		})
	})
}

func TestResponseHeaders_EmptyKeysAndValues(t *testing.T) {
	c, w := newTestContext()
	resp := OK(testPayloadData).
		WithHeader("", "ignored").
		WithHeaderInt("", testEdgeNum100).
		WithHeaderStringer("", &dummyPointerStringer{val: "ignored"}).
		WithHeader("X-Empty-Value", "").
		WithHeaderInt("X-Zero", 0).
		WithHeaderInt("X-Negative", testEdgeNumNeg)

	resp.Write(c)
	assert.NotContains(t, w.Header(), "")
	assert.Equal(t, "", w.Header().Get("X-Empty-Value"))
	assert.Equal(t, "0", w.Header().Get("X-Zero"))
	assert.Equal(t, "-50", w.Header().Get("X-Negative"))

	assert.NotPanics(t, func() { resp.Write(nil) })
	assert.NotPanics(t, func() { NoContent().WithHeaderInt("X-K", testEdgeNum10).Write(nil) })
}

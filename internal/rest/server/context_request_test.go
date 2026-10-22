package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testCtxKey string

const (
	ctxKeyA testCtxKey = "keyA"
	ctxKeyB testCtxKey = "keyB"
	ctxValA            = "valA"
	ctxValB            = "valB"
)

func makeContextReq(method, path string, key testCtxKey, val string) *http.Request {
	return httptest.NewRequest(method, path, nil).WithContext(
		context.WithValue(context.Background(), key, val),
	)
}

func TestContext_SetRequestSync(t *testing.T) {
	gc, _ := gin.CreateTestContext(httptest.NewRecorder())
	gc.Request = makeContextReq(http.MethodGet, "/a", ctxKeyA, ctxValA)
	c := newContext(gc)
	assert.Equal(t, ctxValA, c.Context().Value(ctxKeyA))
	reqB := makeContextReq(http.MethodPost, "/b", ctxKeyB, ctxValB)
	c.SetRequest(reqB)
	assert.Equal(t, ctxValB, c.Context().Value(ctxKeyB))
	assert.Nil(t, c.Context().Value(ctxKeyA))
	assert.Same(t, reqB, c.ginCtx.Request)
	assert.Equal(t, ctxValB, c.ginCtx.Request.Context().Value(ctxKeyB))
}

func TestContext_ContextValueAndCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxKeyA, ctxValA))
	gc, _ := gin.CreateTestContext(httptest.NewRecorder())
	gc.Request = httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)
	c := newContext(gc)
	assert.Equal(t, ctxValA, c.Context().Value(ctxKeyA))
	cancel()
	assert.ErrorIs(t, c.Context().Err(), context.Canceled)
}

func TestContext_Precedence(t *testing.T) {
	gc, _ := gin.CreateTestContext(httptest.NewRecorder())
	gc.Request = makeContextReq(http.MethodGet, "/g", ctxKeyA, ctxValA)
	c := &Context{ginCtx: gc}
	assert.Equal(t, ctxValA, c.Context().Value(ctxKeyA))
	c.Request = makeContextReq(http.MethodGet, "/d", ctxKeyB, ctxValB)
	assert.Equal(t, ctxValB, c.Context().Value(ctxKeyB))
	assert.Nil(t, c.Context().Value(ctxKeyA))
}

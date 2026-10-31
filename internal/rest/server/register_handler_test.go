package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
)

const (
	expectedHandlerCount = 3
	testMWPath           = "/test-mw"
	testAbortPath        = "/test-abort"
	testChainPath        = "/test-chain"
)

func TestWrapMiddleware_Executes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	executed := false

	mw := wrapMiddleware(func(contracts.Context) {
		executed = true
	})
	engine.Use(mw)
	engine.GET(testMWPath, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, testMWPath, nil))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, executed)
}

func TestInvokeHandler_AbortedContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	handlerCalled := false

	engine.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	})
	engine.GET(testAbortPath, wrapHandler(func(contracts.Context) contracts.Response {
		handlerCalled = true
		return testNoContent()
	}))

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, testAbortPath, nil))

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, handlerCalled)
}

func buildTestMiddlewareChain(steps *[]string) ([]contracts.Middleware, contracts.Handler) {
	mw1 := func(contracts.Context) { *steps = append(*steps, "mw1") }
	mw2 := func(contracts.Context) { *steps = append(*steps, "mw2") }
	h := func(contracts.Context) contracts.Response {
		*steps = append(*steps, "h")
		return testNoContent()
	}
	return []contracts.Middleware{mw1, mw2}, h
}

func executeChainRequest(handlers []gin.HandlerFunc) *httptest.ResponseRecorder {
	engine := gin.New()
	engine.GET(testChainPath, handlers...)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, testChainPath, nil))
	return w
}

func TestBuildHandlers_ChainsMiddlewareAndHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var steps []string
	mws, h := buildTestMiddlewareChain(&steps)

	handlers := buildHandlers(mws, h)
	assert.Len(t, handlers, expectedHandlerCount)

	w := executeChainRequest(handlers)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, []string{"mw1", "mw2", "h"}, steps)
}

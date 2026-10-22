package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func liveContextRegs(mwVal, hVal *any) []*contracts.RRestAPIRegistration {
	return []*contracts.RRestAPIRegistration{{
		Version: 1,
		Apis: []*contracts.ExportableAPI{{
			Path:   "/live",
			Method: contracts.GET,
			Middleware: []contracts.Middleware{
				func(c contracts.Context) { *mwVal = c.Context().Value(ctxKeyA) },
			},
			Handler: func(c contracts.Context) contracts.Response {
				*hVal = c.Context().Value(ctxKeyA)
				return testNoContent()
			},
		}},
	}}
}

func liveSetRequestRegs(handlerVal *any) []*contracts.RRestAPIRegistration {
	return []*contracts.RRestAPIRegistration{{
		Version: 1,
		Apis: []*contracts.ExportableAPI{{
			Path:   "/live",
			Method: contracts.GET,
			Middleware: []contracts.Middleware{
				func(c contracts.Context) {
					if sc, ok := c.(*Context); ok && sc.Request != nil {
						sc.SetRequest(sc.Request.WithContext(
							context.WithValue(sc.Context(), ctxKeyB, ctxValB),
						))
					}
				},
			},
			Handler: func(c contracts.Context) contracts.Response {
				*handlerVal = c.Context().Value(ctxKeyB)
				return testNoContent()
			},
		}},
	}}
}

func newLiveEngine(t *testing.T, regs []*contracts.RRestAPIRegistration) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	require.NoError(t, RegisterRoutesWithVersion(engine, regs, ""))
	return engine
}

func TestContext_LiveHTTPContext(t *testing.T) {
	var mwVal, hVal any
	engine := newLiveEngine(t, liveContextRegs(&mwVal, &hVal))
	req := makeContextReq(http.MethodGet, "/api/v1/live", ctxKeyA, ctxValA)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, ctxValA, mwVal)
	assert.Equal(t, ctxValA, hVal)
}

func TestContext_LiveMiddlewareSetRequest(t *testing.T) {
	var handlerVal any
	engine := newLiveEngine(t, liveSetRequestRegs(&handlerVal))
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/live", nil))
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, ctxValB, handlerVal)
}

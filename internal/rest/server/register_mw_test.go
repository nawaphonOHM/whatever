package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mwPath     = "/mw"
	mwFullPath = "/v1/mw"
	mwMarkerA  = "A"
	mwMarkerB  = "B"
	mwMarkerH  = "H"
	mwOrderLen = 3
)

// mwOrderRegs builds a registration that records middleware order.
func mwOrderRegs(order *[]string) []*contracts.RRestAPIRegistration {
	return []*contracts.RRestAPIRegistration{{
		Version: 1,
		Apis: []*contracts.ExportableAPI{{
			Path:   mwPath,
			Method: contracts.GET,
			Middleware: []contracts.Middleware{
				func(contracts.Context) { *order = append(*order, mwMarkerA) },
				func(contracts.Context) { *order = append(*order, mwMarkerB) },
			},
			Handler: func(contracts.Context) contracts.Response {
				*order = append(*order, mwMarkerH)
				return testNoContent()
			},
		}},
	}}
}

// TestRegisterRoutes_MiddlewareOrder checks middleware onion order.
func TestRegisterRoutes_MiddlewareOrder(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	var order []string
	require.NoError(t, RegisterRoutesWithVersion(
		engine, mwOrderRegs(&order), "",
	))

	// Act
	w := httptest.NewRecorder()
	engine.ServeHTTP(
		w,
		httptest.NewRequest(http.MethodGet, mwFullPath, nil),
	)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, []string{mwMarkerA, mwMarkerB, mwMarkerH}, order)
	assert.Len(t, order, mwOrderLen)
}

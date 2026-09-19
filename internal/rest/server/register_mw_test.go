package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/pkg/rest"
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
func mwOrderRegs(order *[]string) []*rest.RestAPIRegistration {
	return []*rest.RestAPIRegistration{{
		Version: 1,
		Apis: []*rest.ExportableAPI{{
			Path:   mwPath,
			Method: rest.GET,
			Middleware: []rest.Middleware{
				func(*rest.Context) { *order = append(*order, mwMarkerA) },
				func(*rest.Context) { *order = append(*order, mwMarkerB) },
			},
			Handler: func(*rest.Context) rest.Response {
				*order = append(*order, mwMarkerH)
				return rest.NoContent()
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

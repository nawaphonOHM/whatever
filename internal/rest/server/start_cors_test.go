package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func corsPreflight(t *testing.T, srv *Server) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	srv.Engine.ServeHTTP(w, req)
	return w
}

func newCORSBlueprint(cors *contracts.CorsSetting) *contracts.BluePrint {
	return contracts.NewBluePrint().WithMeta(contracts.NewMeta().WithCors(cors))
}

// TestNewFromBluePrint_DefaultCORS uses permissive defaults when unset.
func TestNewFromBluePrint_DefaultCORS(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)

	srv, err := NewFromBluePrint(contracts.NewBluePrint())
	require.NoError(t, err)
	response := corsPreflight(t, srv)

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Equal(t, "*", response.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(
		t,
		response.Header().Get("Access-Control-Allow-Methods"),
		http.MethodGet,
	)
}

// TestNewFromBluePrint_CustomCORS maps origins and methods from metadata.
func TestNewFromBluePrint_CustomCORS(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	cors := contracts.NewCorsSetting().
		WithAllowOrigin("https://example.com").
		WithAllowHTTPMethods(contracts.POST, contracts.PATCH)

	srv, err := NewFromBluePrint(newCORSBlueprint(cors))
	require.NoError(t, err)
	response := corsPreflight(t, srv)

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Equal(
		t,
		"https://example.com",
		response.Header().Get("Access-Control-Allow-Origin"),
	)
	assert.Equal(
		t,
		"POST, PATCH",
		response.Header().Get("Access-Control-Allow-Methods"),
	)
}

// TestNewFromBluePrint_EmptyCORS falls back when options are empty.
func TestNewFromBluePrint_EmptyCORS(t *testing.T) {
	t.Setenv(envGinMode, gin.TestMode)
	cors := contracts.NewCorsSetting()
	srv, err := NewFromBluePrint(newCORSBlueprint(cors))
	require.NoError(t, err)
	response := corsPreflight(t, srv)

	assert.Equal(t, "*", response.Header().Get("Access-Control-Allow-Origin"))
	assert.True(t, strings.Contains(
		response.Header().Get("Access-Control-Allow-Methods"),
		http.MethodOptions,
	))
}

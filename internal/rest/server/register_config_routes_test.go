package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRawGroupPath = "/custom/items"
	rawGroupMode     = "raw"
)

func rawRegs() []*contracts.RRestAPIRegistration {
	return []*contracts.RRestAPIRegistration{{
		Group:   contracts.APIGroupRaw,
		Prefix:  "/custom",
		Version: 2,
		Apis: []*contracts.ExportableAPI{{
			Path:   "/items",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK([]string{"raw1", "raw2"})
			},
		}},
	}}
}

func stdRegs() []*contracts.RRestAPIRegistration {
	return []*contracts.RRestAPIRegistration{{
		Group:   contracts.APIGroupStandard,
		Prefix:  "/std",
		Version: 1,
		Apis: []*contracts.ExportableAPI{{
			Path:   "/items",
			Method: contracts.GET,
			Handler: func(contracts.Context) contracts.Response {
				return testOK([]string{"std1"})
			},
		}},
	}}
}

// assertHealthAndReady checks status OK on reserved health endpoints.
func assertHealthAndReady(t *testing.T, engine *gin.Engine) {
	t.Helper()
	code, _ := getJSON(t, engine, ReservedHealthPath)
	assert.Equal(t, http.StatusOK, code)
	code, _ = getJSON(t, engine, ReservedReadyPath)
	assert.Equal(t, http.StatusOK, code)
}

// TestRegisterRoutesWithConfig_Raw mounts raw routes without /api prefix and serves health.
func TestRegisterRoutesWithConfig_Raw(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	cfg := &Config{AppVersion: "2.0.0"}
	require.NoError(t, RegisterRoutesWithConfig(engine, rawRegs(), cfg))

	code, body := getJSON(t, engine, testRawGroupPath)
	assert.Equal(t, http.StatusOK, code)
	assert.Contains(t, string(body), "raw1")
	assertHealthAndReady(t, engine)
}

// TestRegisterRoutesWithConfig_NilConfig derives group from registrations with nil config.
func TestRegisterRoutesWithConfig_NilConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("nil cfg mounts standard routes", func(t *testing.T) {
		engine := gin.New()
		require.NoError(t, RegisterRoutesWithConfig(engine, stdRegs(), nil))
		code, _ := getJSON(t, engine, "/api/v1/std/items")
		assert.Equal(t, http.StatusOK, code)
		assertHealthAndReady(t, engine)
	})

	t.Run("nil cfg mounts raw routes", func(t *testing.T) {
		engine := gin.New()
		require.NoError(t, RegisterRoutesWithConfig(engine, rawRegs(), nil))
		code, body := getJSON(t, engine, testRawGroupPath)
		assert.Equal(t, http.StatusOK, code)
		assert.Contains(t, string(body), "raw1")
		assertHealthAndReady(t, engine)
	})
}

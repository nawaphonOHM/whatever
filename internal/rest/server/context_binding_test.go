package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bindingPayload struct {
	Name string `json:"name" form:"name" uri:"name" header:"X-Name"`
	Age  int    `json:"age" form:"age" uri:"age"`
}

func TestContext_BindJSON(t *testing.T) {
	r := gin.New()
	var payload bindingPayload
	r.POST("/json", func(gc *gin.Context) {
		require.NoError(t, newContext(gc).ShouldBindJSON(&payload))
		gc.Status(testStatus)
	})
	req := httptest.NewRequest(http.MethodPost, "/json",
		strings.NewReader(`{"name":"Alice","age":30}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, testStatus, w.Code)
	assert.Equal(t, "Alice", payload.Name)
	assert.Equal(t, testAgeAlice, payload.Age)
}

func TestContext_BindQuery(t *testing.T) {
	r := gin.New()
	var payload bindingPayload
	r.GET("/query", func(gc *gin.Context) {
		c := newContext(gc)
		require.NoError(t, c.ShouldBindQuery(&payload))
		require.NoError(t, c.BindQuery(&payload))
		gc.Status(testStatus)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet,
		"/query?name=Bob&age=22", nil))
	assert.Equal(t, testStatus, w.Code)
	assert.Equal(t, "Bob", payload.Name)
	assert.Equal(t, testAgeBob, payload.Age)
}

func registerURIAndHeaderRoutes(
	t *testing.T,
	r *gin.Engine,
	uriPayload, headerPayload *bindingPayload,
) {
	t.Helper()
	r.GET("/uri/:name/:age", func(gc *gin.Context) {
		c := newContext(gc)
		require.NoError(t, c.ShouldBindURI(uriPayload))
		require.NoError(t, c.BindURI(uriPayload))
		require.NoError(t, c.ShouldBindHeader(headerPayload))
		require.NoError(t, c.BindHeader(headerPayload))
		gc.Status(testStatus)
	})
}

func TestContext_BindURIAndHeader(t *testing.T) {
	r := gin.New()
	var uriPayload, headerPayload bindingPayload
	registerURIAndHeaderRoutes(t, r, &uriPayload, &headerPayload)
	req := httptest.NewRequest(http.MethodGet, "/uri/Cara/35", nil)
	req.Header.Set("X-Name", "Cara")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, "Cara", uriPayload.Name)
	assert.Equal(t, testAgeCara, uriPayload.Age)
	assert.Equal(t, "Cara", headerPayload.Name)
}

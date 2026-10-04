package rest

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	testValue         = "value"
	invalidMethod     = 999
	testStatus        = 200
	testUnknownStatus = 999
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

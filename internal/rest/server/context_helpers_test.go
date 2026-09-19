package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	testKey          = "key"
	testValue        = "value"
	testDefault      = "default"
	testStatus       = 200
	testBadStatus    = 400
	testAccepted     = 202
	testUnauthorized = 401
	testAgeAlice     = 30
	testAgeBob       = 22
	testAgeCara      = 35
	testIntValue     = 42
	testInt64Value   = 999
	testUintValue    = 12
	testUint64Value  = 888
	testFloatValue   = 3.14
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

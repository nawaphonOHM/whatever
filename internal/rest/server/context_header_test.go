package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func TestContext_HeadersAndCookies(t *testing.T) {
	r := gin.New()
	r.GET("/headers", func(gc *gin.Context) {
		c := newContext(gc)
		assert.Equal(t, "Bearer xyz", c.GetHeader("Authorization"))
		cookie, err := c.Cookie("session")
		assert.NoError(t, err)
		assert.Equal(t, "session-val", cookie)
		c.Header("X-Custom", testValue)
		c.SetHeader("X-Alias", testValue)
		c.SetCookie(&http.Cookie{Name: "remember", Value: testValue})
		c.Status(testAccepted)
	})
	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("Authorization", "Bearer xyz")
	req.AddCookie(&http.Cookie{Name: "session", Value: "session-val"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, testAccepted, w.Code)
	assert.Equal(t, testValue, w.Header().Get("X-Custom"))
	assert.Contains(t, w.Header().Get("Set-Cookie"), "remember=value")
}

func TestContext_ParamsQueriesAndForms(t *testing.T) {
	r := gin.New()
	r.POST("/users/:id", func(gc *gin.Context) {
		c := newContext(gc)
		assert.Equal(t, "42", c.Param("id"))
		assert.Equal(t, "asc", c.Query("sort"))
		assert.Equal(t, testDefault, c.DefaultQuery("missing", testDefault))
		assert.Equal(t, []string{"1", "2"}, c.QueryArray("filter"))
		assert.Equal(t, map[string]string{"a": "1"}, c.QueryMap("tags"))
		assert.Equal(t, "formVal", c.PostForm("field"))
		assert.Equal(t, []string{"x", "y"}, c.PostFormArray("items"))
		props := map[string]string{testKey: testValue}
		assert.Equal(t, props, c.PostFormMap("props"))
		gc.Status(testStatus)
	})
	path := "/users/42?sort=asc&filter=1&filter=2&tags[a]=1"
	body := "field=formVal&items=x&items=y&props[key]=value"
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, testStatus, w.Code)
}

func TestContext_DefaultForms(t *testing.T) {
	gc, _ := newTestContext()
	c := newContext(gc)
	assert.Equal(t, testDefault, c.DefaultPostForm(testKey, testDefault))
	assert.Empty(t, c.PostForm(testKey))
	assert.Nil(t, c.PostFormArray(testKey))
	assert.Empty(t, c.PostFormMap(testKey))
}

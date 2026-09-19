package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

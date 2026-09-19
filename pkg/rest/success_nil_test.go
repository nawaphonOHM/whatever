package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse_NilAndDefaults(t *testing.T) {
	var nilResp *successResponse[string]
	assert.Equal(t, http.StatusOK, nilResp.StatusCode())
	assert.Nil(t, nilResp.withMessage("msg"))
	assert.Nil(t, nilResp.withTimestamp(time.Now()))
	nilResp.Write(nil)
}

func TestSuccessResponse_ZeroCodeWrite(t *testing.T) {
	resp := &successResponse[string]{Code: 0}
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	gc, _ := gin.CreateTestContext(w)
	resp.Write(gc)
	assert.Equal(t, http.StatusOK, w.Code)
}

package rest

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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

func occupyTCPPort(t *testing.T) (net.Listener, int) {
	t.Helper()
	l, err := net.Listen("tcp", testHost+":0")
	require.NoError(t, err)
	tcpAddr, ok := l.Addr().(*net.TCPAddr)
	require.True(t, ok)
	return l, tcpAddr.Port
}

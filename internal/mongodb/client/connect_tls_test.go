package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testErrMsgNonTLS = "connection refused by peer"
	testErrMsgTLS    = "server requires TLS connection"
	testErrMsgSSL    = "SSL handshake failed"
	testErrMsgClosed = "connection closed"
)

// Test_isTLSError tests TLS requirement error pattern matching.
func Test_isTLSError(t *testing.T) {
	assert.False(t, isTLSError(nil))
	assert.False(t, isTLSError(errors.New(testErrMsgNonTLS)))
	assert.False(t, isTLSError(errors.New("i/o timeout")))
	assert.True(t, isTLSError(errors.New(testErrMsgTLS)))
	assert.True(t, isTLSError(errors.New(testErrMsgSSL)))
	assert.True(t, isTLSError(errors.New(testErrMsgClosed)))
	assert.True(t, isTLSError(errors.New("tls: first record does not look like a tls handshake")))
	assert.True(t, isTLSError(errors.New("command requires ssl")))
	assert.True(t, isTLSError(errors.New("requires tls")))
	assert.True(t, isTLSError(errors.New("ssl required")))
}

// Test_isTLSError_NetworkAndResetPatterns tests network reset patterns that indicate TLS requirements.
func Test_isTLSError_NetworkAndResetPatterns(t *testing.T) {
	assert.True(t, isTLSError(errors.New("connection reset by peer")))
	assert.True(t, isTLSError(errors.New("incomplete read of full message")))
	assert.True(t, isTLSError(errors.New("broken pipe")))
	assert.True(t, isTLSError(errors.New("server selection error")))
	assert.True(t, isTLSError(errors.New("read: connection reset by peer")))
	assert.True(t, isTLSError(errors.New("incomplete read of full message: read tcp: connection reset by peer")))
	assert.True(t, isTLSError(errors.New("server selection error: context deadline exceeded")))
	assert.True(t, isTLSError(errors.New("write: broken pipe")))
}

// Test_handleConnectionError tests connection error logging and exit hook execution.
func Test_handleConnectionError(t *testing.T) {
	assert.Nil(t, handleConnectionError(context.Background(), nil))

	var exitCode int
	restore := SetExitFunc(func(code int) {
		exitCode = code
	})
	defer func() { SetExitFunc(restore) }()

	err := errors.New("test connection error")
	returnedErr := handleConnectionError(context.Background(), err)
	assert.Equal(t, err, returnedErr)
	assert.Equal(t, 1, exitCode)
}

// Test_SetExitFunc tests swapping the package-level exitFunc.
func Test_SetExitFunc(t *testing.T) {
	var customCalled bool
	customFunc := func(int) {
		customCalled = true
	}

	prev := SetExitFunc(customFunc)
	defer func() { SetExitFunc(prev) }()
	err := handleConnectionError(context.Background(), errors.New("test exit err"))
	assert.Error(t, err)
	assert.True(t, customCalled)
}

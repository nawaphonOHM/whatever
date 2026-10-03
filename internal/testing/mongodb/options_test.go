package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	require.NotNil(t, opts)
	assert.Equal(t, DefaultHost, opts.Host)
	assert.Equal(t, DefaultPort, opts.Port)
	assert.Equal(t, DefaultProtocol, opts.Protocol)
	assert.Equal(t, DefaultUUIDRepresentation, opts.UUIDRepresentation)
	assert.Equal(t, DefaultConnectTimeout, opts.ConnectTimeout)
	assert.True(t, opts.EnablePing)
	assert.False(t, opts.EnableTLS)
	assert.False(t, opts.DirectConnection)
}

func TestNewOptions_NilSafe(t *testing.T) {
	opts := NewOptions(nil, nil)
	require.NotNil(t, opts)
	assert.Equal(t, DefaultHost, opts.Host)
	assert.Equal(t, DefaultPort, opts.Port)
	assert.Equal(t, DefaultConnectTimeout, opts.ConnectTimeout)
}

const testCustomPort = 27018

func TestWithOptions_Endpoints(t *testing.T) {
	opts := NewOptions(
		WithHost("mongo.test.local"),
		WithPort(testCustomPort),
		WithProtocol("mongodb+srv"),
		WithDatabase("integration_test"),
	)
	assert.Equal(t, "mongo.test.local", opts.Host)
	assert.Equal(t, testCustomPort, opts.Port)
	assert.Equal(t, "mongodb+srv", opts.Protocol)
	assert.Equal(t, "integration_test", opts.Database)
}

func TestWithOptions_AuthAndApp(t *testing.T) {
	opts := NewOptions(
		WithUsername("testuser"),
		WithPassword("testpass"),
		WithAuthSource("admin"),
		WithAppName("test-suite"),
		WithURI("mongodb://localhost:27017/test"),
	)
	assert.Equal(t, "testuser", opts.Username)
	assert.Equal(t, "testpass", opts.Password)
	assert.Equal(t, "admin", opts.AuthSource)
	assert.Equal(t, "test-suite", opts.AppName)
	assert.Equal(t, "mongodb://localhost:27017/test", opts.URI)
}

func TestErrors(t *testing.T) {
	assert.NotNil(t, ErrNilClient)
	assert.NotNil(t, ErrNilConfig)
	assert.NotNil(t, ErrEmptyURI)
	assert.NotNil(t, ErrInvalidPort)
	assert.NotNil(t, ErrNilTestingTB)
}

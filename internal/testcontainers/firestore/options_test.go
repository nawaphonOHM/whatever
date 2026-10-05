package firestore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type dummyCustomizer struct{}

func (*dummyCustomizer) Customize(*testcontainers.GenericContainerRequest) error {
	return nil
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	require.NotNil(t, opts)
	assert.Equal(t, DefaultImage, opts.Image)
	assert.Equal(t, DefaultProjectID, opts.ProjectID)
	assert.False(t, opts.DatastoreMode)
	assert.NotNil(t, opts.Env)
	assert.Empty(t, opts.ContainerOptions)
}

func TestNewOptions_DefaultsAndNil(t *testing.T) {
	opts := NewOptions(nil, nil)
	require.NotNil(t, opts)
	assert.Equal(t, DefaultImage, opts.Image)
	assert.Equal(t, DefaultProjectID, opts.ProjectID)
	assert.False(t, opts.DatastoreMode)
	assert.NotNil(t, opts.Env)
}

func TestWithImage(t *testing.T) {
	opts := NewOptions(WithImage("custom-firestore:latest"))
	assert.Equal(t, "custom-firestore:latest", opts.Image)
}

func TestWithProjectID(t *testing.T) {
	opts := NewOptions(WithProjectID("my-project"))
	assert.Equal(t, "my-project", opts.ProjectID)
}

func TestWithDatastoreMode(t *testing.T) {
	opts1 := NewOptions(WithDatastoreMode())
	assert.True(t, opts1.DatastoreMode)

	opts2 := NewOptions(WithDatastoreMode(true))
	assert.True(t, opts2.DatastoreMode)

	opts3 := NewOptions(WithDatastoreMode(false))
	assert.False(t, opts3.DatastoreMode)
}

func TestWithEnv(t *testing.T) {
	opts := NewOptions(WithEnv("FOO", "BAR"), WithEnv("BAZ", "QUX"))
	assert.Equal(t, "BAR", opts.Env["FOO"])
	assert.Equal(t, "QUX", opts.Env["BAZ"])
}

func TestWithEnv_NilEnv(t *testing.T) {
	opts := &Options{}
	opt := WithEnv("KEY", "VAL")
	opt(opts)
	assert.Equal(t, "VAL", opts.Env["KEY"])
}

func TestWithContainerOptions(t *testing.T) {
	c1 := &dummyCustomizer{}
	opts := NewOptions(WithContainerOptions(c1, nil))
	require.Len(t, opts.ContainerOptions, 1)
	assert.Equal(t, c1, opts.ContainerOptions[0])
}

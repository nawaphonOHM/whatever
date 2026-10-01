package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

type dummyCustomizer struct{}

func (dummyCustomizer) Customize(*testcontainers.GenericContainerRequest) error {
	return nil
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	require.NotNil(t, opts)
	assert.Equal(t, DefaultImage, opts.Image)
	assert.Empty(t, opts.Username)
	assert.Empty(t, opts.Password)
	assert.Empty(t, opts.Database)
	assert.Empty(t, opts.ReplicaSet)
	assert.NotNil(t, opts.Env)
	assert.Empty(t, opts.ContainerOptions)
}

func TestNewOptions_DefaultsAndNil(t *testing.T) {
	opts := NewOptions(nil, nil)
	require.NotNil(t, opts)
	assert.Equal(t, DefaultImage, opts.Image)
	assert.NotNil(t, opts.Env)
}

func TestWithImage(t *testing.T) {
	opts := NewOptions(WithImage("mongo:7.0"))
	assert.Equal(t, "mongo:7.0", opts.Image)
}

func TestWithCredentials(t *testing.T) {
	opts := NewOptions(WithUsername("admin"), WithPassword("secret"))
	assert.Equal(t, "admin", opts.Username)
	assert.Equal(t, "secret", opts.Password)
}

func TestWithDatabaseAndReplicaSet(t *testing.T) {
	opts := NewOptions(WithDatabase("testdb"), WithReplicaSet("rs0"))
	assert.Equal(t, "testdb", opts.Database)
	assert.Equal(t, "rs0", opts.ReplicaSet)
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
	c1 := dummyCustomizer{}
	opts := NewOptions(WithContainerOptions(c1, nil))
	require.Len(t, opts.ContainerOptions, 1)
	assert.Equal(t, c1, opts.ContainerOptions[0])
}

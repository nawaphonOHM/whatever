package azurite

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
	opts := NewOptions(WithImage("custom-azurite:latest"))
	assert.Equal(t, "custom-azurite:latest", opts.Image)
}

func TestWithEnv(t *testing.T) {
	opts := NewOptions(WithEnv("FOO", "BAR"), WithEnv("BAZ", "QUX"))
	assert.Equal(t, "BAR", opts.Env["FOO"])
	assert.Equal(t, "QUX", opts.Env["BAZ"])
}

func TestWithEnv_NilEnv(t *testing.T) {
	opts := &Options{}
	WithEnv("KEY", "VALUE")(opts)
	assert.Equal(t, "VALUE", opts.Env["KEY"])
}

func TestWithContainerOptionsAndBuildCustomizers(t *testing.T) {
	customizer := &dummyCustomizer{}
	opts := NewOptions(WithEnv("KEY", "VALUE"), WithContainerOptions(customizer, nil))
	customizers := buildCustomizers(opts)
	require.Len(t, customizers, 2)
	assert.Same(t, customizer, customizers[1])
}

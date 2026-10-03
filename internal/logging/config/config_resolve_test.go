package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveOutput_Standard(t *testing.T) {
	t.Run("stdout", func(t *testing.T) {
		cfg := &Config{Output: "stdout"}
		w, err := cfg.ResolveOutput()
		require.NoError(t, err)
		assert.Equal(t, os.Stdout, w)
	})

	t.Run("empty defaults to stdout", func(t *testing.T) {
		cfg := &Config{Output: ""}
		w, err := cfg.ResolveOutput()
		require.NoError(t, err)
		assert.Equal(t, os.Stdout, w)
	})

	t.Run(testStderr, func(t *testing.T) {
		cfg := &Config{Output: testStderr}
		w, err := cfg.ResolveOutput()
		require.NoError(t, err)
		assert.Equal(t, os.Stderr, w)
	})

	t.Run(testDiscard, func(t *testing.T) {
		cfg := &Config{Output: testDiscard}
		w, err := cfg.ResolveOutput()
		require.NoError(t, err)
		assert.Equal(t, io.Discard, w)
	})
}

func writeOutputFile(t *testing.T, w io.Writer) {
	t.Helper()
	if closer, ok := w.(io.Closer); ok {
		defer func() {
			require.NoError(t, closer.Close())
		}()
	}
	_, err := io.WriteString(w, "test log message\n")
	assert.NoError(t, err)
}

func testResolveOutputFile(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "app.log")
	cfg := &Config{Output: logPath}
	w, err := cfg.ResolveOutput()
	require.NoError(t, err)
	require.NotNil(t, w)
	writeOutputFile(t, w)

	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "test log message\n", string(content))
}

func TestResolveOutput_FileAndErrors(t *testing.T) {
	t.Run("file writer", func(t *testing.T) {
		testResolveOutputFile(t)
	})

	t.Run("uncreatable file writer error", func(t *testing.T) {
		cfg := &Config{Output: "/nonexistent_dir_99999/app.log"}
		w, err := cfg.ResolveOutput()
		require.ErrorIs(t, err, ErrInvalidOutput)
		assert.Nil(t, w)
	})

	t.Run("nil config", func(t *testing.T) {
		var cfg *Config
		w, err := cfg.ResolveOutput()
		require.ErrorIs(t, err, ErrNilConfig)
		assert.Nil(t, w)
	})
}

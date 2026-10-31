// Package config provides tests for configuration loading logging.
package config

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load_IsSilent(t *testing.T) {
	buf := new(bytes.Buffer)
	orig := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, nil)))
	defer slog.SetDefault(orig)

	cfg, err := Load[TestConfig]()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Empty(t, buf.String(), "rest/config.Load must not emit logs directly to avoid pre-init text leaks")
}

package core_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

func TestBuildHandlerGCPMode(t *testing.T) {
	t.Setenv("OHM9996_CLOUD_PROVIDER", "gcp")
	t.Setenv("OHM9996_OTEL_ENABLED", "false")

	cfg := &config.Config{
		CloudProvider: "gcp",
	}

	var buf bytes.Buffer
	handler := core.BuildHandler(cfg, &buf, nil)
	_, ok := handler.(*core.GCPHandler)
	require.True(t, ok)
}

func TestBuildHandlerStandardMode(t *testing.T) {
	t.Setenv("OHM9996_CLOUD_PROVIDER", "")
	t.Setenv("OHM9996_OTEL_ENABLED", "false")

	cfg := &config.Config{
		CloudProvider:           "",
		DisableTraceCorrelation: true,
	}

	var buf bytes.Buffer
	handler := core.BuildHandler(cfg, &buf, nil)
	_, ok := handler.(*core.GCPHandler)
	assert.False(t, ok)
}

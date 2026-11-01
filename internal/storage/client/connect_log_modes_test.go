package client

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

type storageLogMode struct {
	name     string
	levelKey string
	format   string
}

const (
	formatJSON = "json"
	formatGCP  = "gcp"
	keyLevel   = "level"
	keySev     = "severity"
)

func storageLogModes() []storageLogMode {
	return []storageLogMode{
		{name: formatJSON, levelKey: keyLevel, format: formatJSON},
		{name: formatGCP, levelKey: keySev, format: formatGCP},
	}
}

func (m storageLogMode) handler(output io.Writer) slog.Handler {
	if m.format == formatGCP {
		return core.NewGCPHandler(output, nil)
	}
	return slog.NewJSONHandler(output, nil)
}

func (m storageLogMode) debugHandler(output io.Writer) slog.Handler {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	if m.format == formatGCP {
		return core.NewGCPHandler(output, opts)
	}
	return slog.NewJSONHandler(output, opts)
}

func captureStorageDebugLogs(
	t *testing.T,
	mode storageLogMode,
	cfg *config.Config,
	opts ...Option,
) []map[string]any {
	t.Helper()
	isolateAWSEnv(t)
	var buf bytes.Buffer
	setStorageTestLogger(t, mode.debugHandler(&buf))
	allOpts := append([]Option{WithSkipProbe(true)}, opts...)
	_, err := ConnectWithConfig(context.Background(), cfg, allOpts...)
	require.NoError(t, err)
	return parseStorageLogs(t, buf.Bytes())
}

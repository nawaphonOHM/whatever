package core_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/config"
	"github.com/nawaphonOHM/whatever/v3/internal/logging/core"
)

const standardMixedMessage = `payload: {"action":"create"}`

func TestStandardHandlersPreserveMixedMessage(t *testing.T) {
	for _, format := range []config.Format{config.FormatJSON, config.FormatText} {
		t.Run(string(format), func(t *testing.T) {
			var output bytes.Buffer
			cfg := config.DefaultConfig()
			cfg.Format = string(format)
			cfg.DisableTraceCorrelation = true
			slog.New(core.BuildHandler(cfg, &output, nil)).Info(standardMixedMessage)
			assertStandardMessage(t, format, output.Bytes())
		})
	}
}

func assertStandardMessage(t *testing.T, format config.Format, output []byte) {
	t.Helper()
	if format == config.FormatJSON {
		assertJSONMessageUnchanged(t, output)
		return
	}
	assert.Contains(t, string(output), "msg="+strconv.Quote(standardMixedMessage))
}

func assertJSONMessageUnchanged(t *testing.T, output []byte) {
	t.Helper()
	var record map[string]any
	require.NoError(t, json.Unmarshal(output, &record))
	assert.Equal(t, standardMixedMessage, record["msg"])
}

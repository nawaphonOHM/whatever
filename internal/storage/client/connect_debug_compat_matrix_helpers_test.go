package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

type debugCompatExpected struct {
	provider           string
	reqChecksumCalc    string
	respChecksumVal    string
	checksumBehavior   string
	pathStyleBehavior  string
	authMode           string
	usePathStyle       bool
	disableSkipLog     bool
	hasDriverOverrides bool
}

type debugCompatCase struct {
	configure func(*config.Config)
	name      string
	opts      []Option
	expected  debugCompatExpected
}

func assertCompatChecksumFields(t *testing.T, entry map[string]any, exp debugCompatExpected) {
	t.Helper()
	assert.Equal(t, exp.provider, entry["provider"])
	assert.Equal(t, exp.reqChecksumCalc, entry["request_checksum_calculation"])
	assert.Equal(t, exp.respChecksumVal, entry["response_checksum_validation"])
	assert.Equal(t, exp.checksumBehavior, entry["checksum_behavior_change"])
}

func assertCompatPathAndAuthFields(t *testing.T, entry map[string]any, exp debugCompatExpected) {
	t.Helper()
	assert.Equal(t, exp.usePathStyle, entry["use_path_style"])
	assert.Equal(t, exp.pathStyleBehavior, entry["path_style_behavior_change"])
	assert.Equal(t, exp.disableSkipLog, entry["disable_log_output_checksum_validation_skipped"])
	assert.Equal(t, exp.authMode, entry["auth_mode"])
	assert.Equal(t, exp.hasDriverOverrides, entry["has_driver_overrides"])
}

func assertCompatLogFields(
	t *testing.T,
	entry map[string]any,
	exp debugCompatExpected,
	mode storageLogMode,
) {
	t.Helper()
	require.NotNil(t, entry, "compat debug log must be present")
	assert.Equal(t, "DEBUG", entry[mode.levelKey])
	assertCompatChecksumFields(t, entry, exp)
	assertCompatPathAndAuthFields(t, entry, exp)
}

func assertCompatCaseLogs(t *testing.T, tc debugCompatCase, mode storageLogMode) {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	tc.configure(cfg)
	entries := captureStorageDebugLogs(t, mode, cfg, tc.opts...)
	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	assertCompatLogFields(t, compatLog, tc.expected, mode)
}

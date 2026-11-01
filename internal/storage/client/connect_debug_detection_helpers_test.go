package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

type debugDetectionExpected struct {
	rawEndpoint        string
	rawBucketURL       string
	inspectedHost      string
	detectedProvider   string
	detectionRule      string
	extractedBucket    string
	normalizedEndpoint string
}

type debugDetectionCase struct {
	name      string
	configure func(*config.Config)
	expected  debugDetectionExpected
}

func assertDetectionLogFields(
	t *testing.T,
	entry map[string]any,
	exp debugDetectionExpected,
	mode storageLogMode,
) {
	t.Helper()
	require.NotNil(t, entry, "detection debug log must be present")
	assert.Equal(t, "DEBUG", entry[mode.levelKey])
	assert.Equal(t, exp.rawEndpoint, entry["raw_endpoint"])
	assert.Equal(t, exp.rawBucketURL, entry["raw_bucket_url"])
	assert.Equal(t, exp.inspectedHost, entry["inspected_host"])
	assert.Equal(t, exp.detectedProvider, entry["detected_provider"])
	assert.Equal(t, exp.detectionRule, entry["detection_rule"])
	assert.Equal(t, exp.extractedBucket, entry["extracted_bucket"])
	assert.Equal(t, exp.normalizedEndpoint, entry["normalized_endpoint"])
}

func assertDetectionCaseLogs(t *testing.T, tc debugDetectionCase, mode storageLogMode) {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	tc.configure(cfg)
	entries := captureStorageDebugLogs(t, mode, cfg)
	detectLog := findStorageLog(entries, "evaluating storage provider detection")
	assertDetectionLogFields(t, detectLog, tc.expected, mode)
}

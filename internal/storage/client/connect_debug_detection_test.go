package client

import (
	"testing"
)

func runDetectionCaseList(t *testing.T, cases []debugDetectionCase, mode storageLogMode) {
	t.Helper()
	for _, tc := range cases {
		t.Run(mode.name+"_"+tc.name, func(t *testing.T) {
			assertDetectionCaseLogs(t, tc, mode)
		})
	}
}

func TestConnect_DebugDetection_AWS(t *testing.T) {
	for _, mode := range storageLogModes() {
		runDetectionCaseList(t, detectionCasesAWS(), mode)
	}
}

func TestConnect_DebugDetection_NonAWS(t *testing.T) {
	for _, mode := range storageLogModes() {
		runDetectionCaseList(t, detectionCasesNonAWS(), mode)
	}
}

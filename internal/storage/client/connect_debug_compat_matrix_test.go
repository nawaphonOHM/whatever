package client

import (
	"testing"
)

func runCompatCaseList(t *testing.T, cases []debugCompatCase, mode storageLogMode) {
	t.Helper()
	for _, tc := range cases {
		t.Run(mode.name+"_"+tc.name, func(t *testing.T) {
			assertCompatCaseLogs(t, tc, mode)
		})
	}
}

func TestConnect_DebugCompat_GCSAndS3(t *testing.T) {
	for _, mode := range storageLogModes() {
		runCompatCaseList(t, compatMatrixCasesGCSAndS3(), mode)
	}
}

func TestConnect_DebugCompat_AWS(t *testing.T) {
	for _, mode := range storageLogModes() {
		runCompatCaseList(t, compatMatrixCasesAWS(), mode)
	}
}

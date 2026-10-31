package core_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isBasicEmptyOrPass(s string) bool {
	if s == "" {
		return true
	}
	return s == "PASS"
}

func isHarnessPrefix(s string) bool {
	if strings.HasPrefix(s, "===") {
		return true
	}
	return strings.HasPrefix(s, "---")
}

func isCoverageReportLine(s string) bool {
	if s == "coverage: [no statements]" {
		return true
	}
	return strings.HasPrefix(s, "coverage: ") && strings.Contains(s, "% of statements")
}

func isTestHarnessLine(s string) bool {
	return isBasicEmptyOrPass(s) || isHarnessPrefix(s) || isCoverageReportLine(s)
}

type harnessTestCase struct {
	name     string
	line     string
	expected bool
}

func harnessPositiveTestCases() []harnessTestCase {
	return []harnessTestCase{
		{name: "empty line", line: "", expected: true},
		{name: "pass line", line: "PASS", expected: true},
		{name: "run prefix", line: "=== RUN   TestFoo", expected: true},
		{name: "pass prefix", line: "--- PASS: TestFoo (0.00s)", expected: true},
		{name: "coverage report", line: "coverage: 14.6% of statements", expected: true},
		{name: "coverage zero percent", line: "coverage: 0.0% of statements", expected: true},
		{name: "coverage with coverpkg", line: "coverage: 14.6% of statements in ./...", expected: true},
		{name: "coverage no statements", line: "coverage: [no statements]", expected: true},
	}
}

func harnessNegativeTestCases() []harnessTestCase {
	return []harnessTestCase{
		{name: "plain text log", line: "INFO loading REST configuration", expected: false},
		{name: "plain text error", line: "failed to connect to mongodb", expected: false},
		{name: "plain text with coverage prefix", line: "coverage: computing totals", expected: false},
		{name: "json log line", line: `{"severity":"INFO","message":"test"}`, expected: false},
	}
}

func TestIsTestHarnessLine_Positive(t *testing.T) {
	for _, tt := range harnessPositiveTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isTestHarnessLine(tt.line))
		})
	}
}

func TestIsTestHarnessLine_Negative(t *testing.T) {
	for _, tt := range harnessNegativeTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isTestHarnessLine(tt.line))
		})
	}
}

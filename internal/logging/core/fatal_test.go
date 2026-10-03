// Package core provides internal logging engine implementation.
package core

import (
	"bytes"
	"sync/atomic"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFatalTestLogger(buf *bytes.Buffer, exitCode *atomic.Int32) *Logger {
	l := NewJSON(buf, config.LevelInfo)
	l.SetExitFunc(func(code int) {
		exitCode.Store(int32(code))
	})
	return l
}

func TestLogger_Fatal_CustomExit(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	l := createFatalTestLogger(buf, &exitCode)
	l.Fatal("fatal error occurred", "reason", "db_down")

	assert.Equal(t, int32(1), exitCode.Load())
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "fatal error occurred", logMap["msg"])
	assert.Equal(t, "FATAL", logMap["level"])
}

func TestLogger_FatalContext_WithSpanAndCustomExit(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	l := createFatalTestLogger(buf, &exitCode)
	ctx := createTestSpanContext(t)
	l.FatalContext(ctx, "fatal context error")

	assert.Equal(t, int32(1), exitCode.Load())
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", logMap["trace_id"])
}

func TestLogger_SetExitFunc(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	l := NewJSON(buf, config.LevelInfo)
	l.SetExitFunc(func(code int) { exitCode.Store(int32(code)) })

	require.NotNil(t, l.ExitFunc())
	l.Fatal("fatal with dynamically set exit func")

	assert.Equal(t, int32(1), exitCode.Load())
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "FATAL", logMap["level"])
}

func TestPackageLevel_Fatal(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	l := createFatalTestLogger(buf, &exitCode)
	cleanup := SetTestLogger(l)
	defer cleanup()

	Fatal("package level fatal")
	assert.Equal(t, int32(1), exitCode.Load())
}

func TestPackageLevel_FatalContext(t *testing.T) {
	buf := &bytes.Buffer{}
	var exitCode atomic.Int32
	exitCode.Store(-1)

	l := createFatalTestLogger(buf, &exitCode)
	cleanup := SetTestLogger(l)
	defer cleanup()

	ctx := createTestSpanContext(t)
	FatalContext(ctx, "package level fatal context")
	assert.Equal(t, int32(1), exitCode.Load())
}

package core

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/nawaphonOHM/whatever/internal/logging/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type lifecycleTrackingCloser struct {
	bytes.Buffer
	closeCount int
	isClosed   bool
}

func (tc *lifecycleTrackingCloser) Close() error {
	tc.closeCount++
	tc.isClosed = true
	return nil
}

type dummyErrorCloser struct {
	bytes.Buffer
}

func (*dummyErrorCloser) Close() error {
	return errors.New("err close")
}

func emitGlobalLoggingCalls(ctx context.Context) {
	Trace("msg trace")
	TraceContext(ctx, "msg trace ctx")
	Debug("msg debug")
	Info("msg info")
	Warn("msg warn")
	Error("msg error")
	Log(ctx, slog.LevelInfo, "msg log")
	LogAttrs(ctx, slog.LevelWarn, "msg log attrs", slog.String("k", "v"))
}

func assertCloserState(t *testing.T, tc *lifecycleTrackingCloser, closed bool, count int) {
	t.Helper()
	assert.Equal(t, closed, tc.isClosed)
	assert.Equal(t, count, tc.closeCount)
}

func TestGlobalLogger_SingletonResourceProof(t *testing.T) {
	defer ResetDefault()

	tc := &lifecycleTrackingCloser{}
	l := NewJSON(tc, config.LevelTrace)
	require.NotNil(t, l.closer)
	SetDefault(l)

	emitGlobalLoggingCalls(context.Background())
	assertCloserState(t, tc, false, 0)

	ResetDefault()
	assertCloserState(t, tc, true, 1)
}

func TestClosePreviousLogger_NilAndErrorSafety(t *testing.T) {
	assert.NotPanics(t, func() { closePreviousLogger(nil) })
	assert.NotPanics(t, func() { closePreviousLogger(&Logger{}) })

	errCloser := &dummyErrorCloser{}
	errLogger := NewJSON(errCloser, config.LevelInfo)
	assert.NotPanics(t, func() { closePreviousLogger(errLogger) })
}

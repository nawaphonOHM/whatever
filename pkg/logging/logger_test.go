package logging

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type levelFilterCase struct {
	logAction     func(*Logger)
	name          string
	expectedLevel string
	cfgLevel      Level
	expectLog     bool
}

func levelFilterCases() []levelFilterCase {
	return []levelFilterCase{
		{
			logAction: func(l *Logger) { l.Trace("trace") },
			name:      "trace suppressed at info",
			cfgLevel:  LevelInfo,
			expectLog: false,
		},
		{
			logAction:     func(l *Logger) { l.Info("info") },
			name:          "info emitted at info",
			expectedLevel: "INFO",
			cfgLevel:      LevelInfo,
			expectLog:     true,
		},
		{
			logAction:     func(l *Logger) { l.Trace("trace") },
			name:          "trace emitted at trace",
			expectedLevel: "TRACE",
			cfgLevel:      LevelTrace,
			expectLog:     true,
		},
		{
			logAction:     func(l *Logger) { l.Warn("warn") },
			name:          "warn emitted at warn",
			expectedLevel: "WARN",
			cfgLevel:      LevelWarn,
			expectLog:     true,
		},
	}
}

func TestNew_Defaults(t *testing.T) {
	buf := &bytes.Buffer{}
	l := New(Config{Output: buf})
	require.NotNil(t, l)

	l.Info("default info message")
	logMap := parseJSONLog(t, buf.Bytes())
	assert.Equal(t, "default info message", logMap["msg"])
	assert.Equal(t, "INFO", logMap["level"])
}

func TestNew_LevelFiltering(t *testing.T) {
	for _, tt := range levelFilterCases() {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := New(Config{Output: buf, Level: tt.cfgLevel})
			tt.logAction(l)

			if !tt.expectLog {
				assert.Empty(t, buf.String())
				return
			}
			logMap := parseJSONLog(t, buf.Bytes())
			assert.Equal(t, tt.expectedLevel, logMap["level"])
		})
	}
}

func TestNewJSON_And_NewText(t *testing.T) {
	t.Run("NewJSON", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := NewJSON(buf, LevelDebug)
		l.Debug("json debug message")
		logMap := parseJSONLog(t, buf.Bytes())
		assert.Equal(t, "DEBUG", logMap["level"])
	})

	t.Run("NewText", func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := NewText(buf, LevelInfo)
		l.Info("text info message")
		assert.Contains(t, buf.String(), "level=INFO")
	})
}

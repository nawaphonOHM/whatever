package logging

import (
	"context"
)

// Trace logs at SlogLevelTrace with the default background context.
func (l *Logger) Trace(msg string, args ...any) {
	l.TraceContext(context.Background(), msg, args...)
}

// TraceContext logs at SlogLevelTrace with the provided context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	if l != nil && l.Logger != nil {
		l.Log(ctx, SlogLevelTrace, msg, args...)
	}
}

// Fatal logs at SlogLevelFatal with the default background context and invokes ExitFunc.
func (l *Logger) Fatal(msg string, args ...any) {
	l.FatalContext(context.Background(), msg, args...)
}

func (l *Logger) isNil() bool {
	return l == nil || l.Logger == nil
}

func (l *Logger) executeExit() {
	if l.exitFunc != nil {
		l.exitFunc(1)
	}
}

// FatalContext logs at SlogLevelFatal with the provided context and invokes ExitFunc.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	if l.isNil() {
		return
	}
	l.Log(ctx, SlogLevelFatal, msg, args...)
	l.executeExit()
}

// With returns a new Logger with the specified attributes added.
func (l *Logger) With(args ...any) *Logger {
	if l == nil || l.Logger == nil {
		return l
	}
	return &Logger{
		Logger:   l.Logger.With(args...),
		exitFunc: l.exitFunc,
		handler:  l.handler,
	}
}

// WithGroup returns a new Logger that starts a group with the given name.
func (l *Logger) WithGroup(name string) *Logger {
	if l == nil || l.Logger == nil {
		return l
	}
	return &Logger{
		Logger:   l.Logger.WithGroup(name),
		exitFunc: l.exitFunc,
		handler:  l.handler,
	}
}

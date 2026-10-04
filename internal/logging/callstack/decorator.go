// Package callstack provides call stack tracking, decorators, and context propagation.
package callstack

import (
	"context"
	"runtime"
	"strings"
	"time"
)

func applyFuncInfo(fn *runtime.Func, frame *Frame) {
	pkg, fnName := splitPackageAndFunc(fn.Name())
	frame.Package = pkg
	if frame.Function == "" {
		frame.Function = fnName
	}
}

func resolveCallerInfo(pc uintptr, frame *Frame) {
	if frame == nil {
		return
	}
	if fn := runtime.FuncForPC(pc); fn != nil {
		applyFuncInfo(fn, frame)
	}
}

func captureCallerFrame(name string, skip int) *Frame {
	frame := &Frame{
		Function:  name,
		Timestamp: time.Now(),
	}
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		frame.File = file
		frame.Line = line
		resolveCallerInfo(pc, frame)
	}
	return frame
}

func splitPackageAndFunc(fullName string) (prefix, rest string) {
	lastSlash := strings.LastIndex(fullName, "/")
	prefix = ""
	rest = fullName
	if lastSlash >= 0 {
		prefix = fullName[:lastSlash+1]
		rest = fullName[lastSlash+1:]
	}
	firstDot := strings.Index(rest, ".")
	if firstDot >= 0 {
		return prefix + rest[:firstDot], rest[firstDot+1:]
	}
	return prefix, rest
}

// Decorate wraps a parameterless func() and tracks its frame on entry and exit.
func Decorate(name string, fn func()) func() {
	return func() {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		fn()
	}
}

// DecorateErr wraps a func() error and tracks its frame on entry and exit.
func DecorateErr(name string, fn func() error) func() error {
	return func() error {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		return fn()
	}
}

// DecorateContext wraps a context-aware function func(context.Context) error.
// It ensures a CallStack is attached to context, pushes the frame on entry and pops on exit.
func DecorateContext(name string, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		ctxWithStack, cs := EnsureContext(ctx)
		frame := captureCallerFrame(name, 2)
		cs.Push(frame)
		defer cs.Pop()
		return fn(ctxWithStack)
	}
}

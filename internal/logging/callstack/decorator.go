// Package callstack provides call stack tracking, decorators, and context propagation.
package callstack

import (
	"context"
)

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

// DecorateFunc wraps a generic func() T.
func DecorateFunc[T any](name string, fn func() T) func() T {
	return func() T {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		return fn()
	}
}

// DecorateFuncErr wraps a generic func() (T, error).
func DecorateFuncErr[T any](name string, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		return fn()
	}
}

// DecorateContextFunc wraps a generic func(context.Context) (T, error).
func DecorateContextFunc[T any](name string, fn func(context.Context) (T, error)) func(context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		ctxWithStack, cs := EnsureContext(ctx)
		frame := captureCallerFrame(name, 2)
		cs.Push(frame)
		defer cs.Pop()
		return fn(ctxWithStack)
	}
}

// DecorateCtx wraps a generic func(context.Context) T.
func DecorateCtx[T any](name string, fn func(context.Context) T) func(context.Context) T {
	return func(ctx context.Context) T {
		ctxWithStack, cs := EnsureContext(ctx)
		frame := captureCallerFrame(name, 2)
		cs.Push(frame)
		defer cs.Pop()
		return fn(ctxWithStack)
	}
}
